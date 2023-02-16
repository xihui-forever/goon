package mem

import (
	"github.com/darabuchi/utils"
	"github.com/xihui-forever/goon/middleware/storage"
	"sync"
	"time"
)

type Mem struct {
	lock sync.RWMutex

	data      map[string]*ItemCount
	autoClean bool
}

type ItemCount struct {
	lock sync.RWMutex

	itemList []interface{}
	value    string
	expireAt time.Time
}

func (p *ItemCount) GetValue() string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.value
}

func (p *ItemCount) SetValue(value int64) {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.value = utils.ToString(value)
}

func (p *ItemCount) GetItemList() []interface{} {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.itemList
}

func (p *ItemCount) SetItemList(itemList []interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.itemList = itemList
}

func (p *ItemCount) GetExpireAt() time.Time {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.expireAt
}

func (p *ItemCount) SetExpireAt(expireAt time.Time) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.expireAt = expireAt
}

func (p *ItemCount) SetExpire(expire time.Duration) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.expireAt = time.Now().Add(expire)
}

func (p *ItemCount) GetValid(t time.Duration) int64 {
	for _, item := range p.itemList {
		if time.Now().Sub(item.(time.Time)) > t {
			p.DecBy(utils.ToInt64(item))
		}
	}
	return int64(len(p.GetValue()))
}

func (p *ItemCount) Expire() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if !p.expireAt.IsZero() && p.expireAt.Before(time.Now()) {
		return true
	}
	return false
}

func (p *ItemCount) IncBy(c int64) int64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	itemList := append(p.itemList, c)
	p.SetItemList(itemList)
	return utils.ToInt64(p.GetValue())
}

func (p *ItemCount) DecBy(value int64) int64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	var itemList []interface{}
	for k, v := range p.itemList {
		if utils.ToInt64(v)-value >= value {
			itemList = append(p.itemList[:k], p.itemList[k+1:]...)
		}
	}
	p.SetItemList(itemList)
	return utils.ToInt64(len(p.itemList))
}

func (p *Mem) Clean(logic func(int64) bool) {
	for key, item := range p.Clone() {
		if item.Expire() {
			p.lock.Lock()
			_, ok := p.data[key]
			if ok {
				delete(p.data, key)
			}
			p.lock.Unlock()

			if ok && logic(int64(len(item.value))) {
				return
			}
		}
	}
}

func (p *Mem) needClean() bool {
	if !p.autoClean {
		return false
	}

	p.lock.RLock()
	defer p.lock.RUnlock()

	return len(p.data) > 100
}

func (p *Mem) Clone() map[string]*ItemCount {
	p.lock.RLock()
	defer p.lock.RUnlock()

	data := make(map[string]*ItemCount)
	for k, v := range p.data {
		data[k] = v
	}

	return data
}

func (p *Mem) SetNx(key string, value string) (bool, error) {
	p.lock.RLock()
	_, ok := p.data[key]
	p.lock.RUnlock()

	if ok {
		return false, nil
	}

	defer func() {
		if p.needClean() {
			p.Clean(func(i int64) bool {
				return true
			})
		}
	}()

	p.lock.Lock()
	defer p.lock.Unlock()
	_, ok = p.data[key]

	if ok {
		return false, nil
	}

	p.data[key] = &ItemCount{
		value: value,
	}

	return true, nil
}

func (p *Mem) Get(key string) (string, error) {
	p.lock.RLock()
	item, ok := p.data[key]
	p.lock.RUnlock()
	if !ok {
		return "", storage.ErrKeyNotExist
	}

	return item.GetValue(), nil
}

func (p *Mem) Inc(key string) (int64, error) {
	return p.IncBy(key, time.Now().Unix())
}

func (p *Mem) IncBy(key string, value int64) (int64, error) {
	defer func() {
		if p.needClean() {
			p.Clean(func(i int64) bool {
				return true
			})
		}
	}()

	p.lock.RLock()
	item, ok := p.data[key]
	p.lock.RUnlock()
	if !ok {
		// 不存在就设置默认值并且重新获取一下
		_, err := p.SetNx(key, "0")
		if err != nil {
			return 0, err
		}

		p.lock.RLock()
		item, ok = p.data[key]
		p.lock.RUnlock()
		if !ok {
			return 0, storage.ErrKeyNotExist
		}
	}

	return item.IncBy(value), nil
}

func (p *Mem) Dec(key string) (int64, error) {
	return p.DecBy(key, time.Now().Unix())
}

func (p *Mem) DecBy(key string, value int64) (int64, error) {
	defer func() {
		if p.needClean() {
			p.Clean(func(i int64) bool {
				return true
			})
		}
	}()

	p.lock.RLock()
	item, ok := p.data[key]
	p.lock.RUnlock()
	if !ok {
		// 不存在就设置默认值并且重新获取一下
		_, err := p.SetNx(key, "0")
		if err != nil {
			return 0, err
		}

		p.lock.RLock()
		item, ok = p.data[key]
		p.lock.RUnlock()
		if !ok {
			return 0, storage.ErrKeyNotExist
		}
	}

	return item.DecBy(value), nil
}

func (p *Mem) Expire(key string, timeout time.Duration) error {
	p.lock.Lock()
	item, ok := p.data[key]
	if !ok {
		return storage.ErrKeyNotExist
	}
	p.lock.Unlock()

	if item.Expire() {
		p.lock.Lock()
		delete(p.data, key)
		p.lock.Unlock()
		return storage.ErrKeyNotExist
	}

	item.SetExpire(timeout)

	return nil
}

func (p *Mem) Close() error {
	return nil
}

func NewMem(autoClean bool) *Mem {
	p := &Mem{
		autoClean: autoClean,
		data:      map[string]*ItemCount{},
	}

	return p
}
