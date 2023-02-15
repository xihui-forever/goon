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

	startAt  time.Time
	count    []int64
	value    string
	expireAt time.Time
}

func (p *ItemCount) GetValue() string {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.value
}

func (p *ItemCount) SetValue() {
	p.lock.Lock()
	defer p.lock.Unlock()
	var val int64
	for _, v := range p.count {
		val = val + v
	}
	p.value = utils.ToString(val)
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
	p.startAt = time.Now()
	p.expireAt = p.startAt.Add(expire)
}

func (p *ItemCount) Expire() bool {
	p.lock.RLock()
	defer p.lock.RUnlock()
	if !p.expireAt.IsZero() && p.expireAt.Equal(time.Now()) {
		p.startAt = p.startAt.Add(10 * time.Second)
		p.expireAt = time.Now().Add(10 * time.Second)
		return true
	}
	return false
}

func (p *ItemCount) IncBy(c int64, k int64) int64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.count[k] = p.count[k] + c
	p.SetValue()
	return utils.ToInt64(p.GetValue())
}

func (p *ItemCount) DecBy(c int64, k int64) int64 {
	p.lock.Lock()
	defer p.lock.Unlock()

	p.count[k] = p.count[k] - c
	p.SetValue()
	return utils.ToInt64(p.GetValue())
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

	if item.Expire() {
		p.lock.Lock()
		p.data[key].count = append(p.data[key].count[:1], p.data[key].count[2:]...)
		p.data[key].SetValue()
		p.lock.Unlock()
	}

	return item.GetValue(), nil
}

func (p *Mem) Inc(key string) (int64, error) {
	return p.IncBy(key, 1)
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

	if item.Expire() {
		p.lock.Lock()
		p.data[key].count = append(p.data[key].count[:1], p.data[key].count[2:]...)
		p.data[key].SetValue()
		p.lock.Unlock()
	}

	d := time.Now().Sub(item.startAt)/(10*time.Second) + 1
	k := utils.ToInt64(d)

	return item.IncBy(value, k), nil
}

func (p *Mem) Dec(key string) (int64, error) {
	return p.DecBy(key, 1)
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

	if item.Expire() {
		p.lock.Lock()
		p.data[key].count = append(p.data[key].count[:1], p.data[key].count[2:]...)
		p.data[key].SetValue()
		p.lock.Unlock()
	}

	d := time.Now().Sub(item.startAt)/(10*time.Second) + 1
	k := utils.ToInt64(d)

	return item.DecBy(value, k), nil
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
		p.data[key].count = append(p.data[key].count[:1], p.data[key].count[2:]...)
		p.data[key].SetValue()
		p.lock.Unlock()
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
