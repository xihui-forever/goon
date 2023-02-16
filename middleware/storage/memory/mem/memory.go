package mem

import (
	"github.com/darabuchi/utils"
	"sync"
	"time"

	"github.com/xihui-forever/goon/middleware/storage"
)

/*
		TODO: 1. 通过协程清理过期的key
	          2. 在用户调用方法时，清理过期的key
	          3. 允许用户选择清理的机制（1或者是2）
*/
type Memory struct {
	lock sync.RWMutex

	data      map[string]*Item
	autoClean bool
}

type Item struct {
	lock sync.RWMutex

	itemList []interface{}
	value    string
	expireAt time.Time
}

func (p *Item) GetItemList() []interface{} {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.itemList
}

func (p *Item) SetItemList(itemList []interface{}) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.itemList = itemList
}

func (p *Item) GetValue() string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.value
}

func (p *Item) SetValue(value string) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.value = value
}

func (p *Item) GetExpireAt() time.Time {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.expireAt
}

func (p *Item) SetExpireAt(expireAt time.Time) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.expireAt = expireAt
}

func (p *Item) AddItem(value int64) {
	p.lock.Lock()
	defer p.lock.Unlock()
	items := append(p.itemList, value)
	p.SetItemList(items)
}

func (p *Item) SetExpire(expire time.Duration) {
	p.lock.Lock()
	defer p.lock.Unlock()
	p.expireAt = time.Now().Add(expire)
}

func (p *Item) Expire() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if !p.expireAt.IsZero() && p.expireAt.Before(time.Now()) {
		return true
	}
	return false
}

func (p *Item) IncBy(c int64) int64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	value := utils.ToInt64(p.value) + c
	p.value = utils.ToString(value)
	return value
}

func (p *Item) DecBy(c int64) int64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	value := utils.ToInt64(p.value) - c
	p.value = utils.ToString(value)
	return value
}

func (p *Item) Len() int64 {
	return int64(len(p.itemList))
}

func (p *Memory) Clean(logic func(int64) bool) {
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

func (p *Memory) needClean() bool {
	if !p.autoClean {
		return false
	}

	p.lock.RLock()
	defer p.lock.RUnlock()

	return len(p.data) > 100
}

func (p *Memory) Clone() map[string]*Item {
	p.lock.RLock()
	defer p.lock.RUnlock()

	data := make(map[string]*Item)
	for k, v := range p.data {
		data[k] = v
	}

	return data
}

func (p *Memory) SetNx(key string, value string) (bool, error) {
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

	p.data[key] = &Item{
		value: value,
	}

	return true, nil
}

func (p *Memory) Get(key string) (string, error) {
	p.lock.RLock()
	item, ok := p.data[key]
	p.lock.RUnlock()
	if !ok {
		return "", storage.ErrKeyNotExist
	}

	if item.Expire() {
		p.lock.Lock()
		delete(p.data, key)
		p.lock.Unlock()
		return "", storage.ErrKeyNotExist
	}

	return item.GetValue(), nil
}

func (p *Memory) Inc(key string) (int64, error) {
	return p.IncBy(key, 1)
}

func (p *Memory) IncBy(key string, value int64) (int64, error) {
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

	if item.Expire() {
		p.lock.Lock()
		delete(p.data, key)
		p.lock.Unlock()
		return 0, storage.ErrKeyNotExist
	}

	return item.IncBy(value), nil
}

func (p *Memory) Dec(key string) (int64, error) {
	return p.DecBy(key, 1)
}

func (p *Memory) DecBy(key string, value int64) (int64, error) {
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

	if item.Expire() {
		p.lock.Lock()
		delete(p.data, key)
		p.lock.Unlock()
		return 0, storage.ErrKeyNotExist
	}

	return item.DecBy(value), nil
}

func (p *Memory) Expire(key string, timeout time.Duration) error {
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

func (p *Memory) AddItem(key string, value int64) error {
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
			return err
		}

		p.lock.RLock()
		item, ok = p.data[key]
		p.lock.RUnlock()
		if !ok {
			return storage.ErrKeyNotExist
		}
	}

	if item.Expire() {
		p.lock.Lock()
		delete(p.data, key)
		p.lock.Unlock()
		return storage.ErrKeyNotExist
	}
	item.AddItem(value)
	return nil
}

func (p *Memory) GetNotValid(key string, timeout time.Duration) (int64, error) {
	p.lock.Lock()
	item, ok := p.data[key]
	if !ok {
		return 0, storage.ErrKeyNotExist
	}
	p.lock.Unlock()

	if item.Expire() {
		p.lock.Lock()
		delete(p.data, key)
		p.lock.Unlock()
		return 0, storage.ErrKeyNotExist
	}

	for key, value := range item.itemList {
		if time.Now().Sub(value.(time.Time)) <= timeout {
			return int64(key), nil
		}
	}
	return 0, nil
}

func (p *Memory) DeleteItem(key string, value int64) error {
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

	items := append(item.itemList[:value-1], item.itemList[value:]...)
	item.SetItemList(items)
	return nil
}

func (p *Memory) LenItemList(key string) (int64, error) {
	p.lock.Lock()
	item, ok := p.data[key]
	if !ok {
		return 0, storage.ErrKeyNotExist
	}
	p.lock.Unlock()

	if item.Expire() {
		p.lock.Lock()
		delete(p.data, key)
		p.lock.Unlock()
		return 0, storage.ErrKeyNotExist
	}
	return item.Len(), nil
}

func (p *Memory) Close() error {
	return nil
}

func New(autoClean bool) *Memory {
	p := &Memory{
		autoClean: autoClean,
		data:      map[string]*Item{},
	}

	return p
}
