package lru

import (
	"sync"
	"time"

	"github.com/elliotchance/pie/v2"
	"github.com/xihui-forever/goon/storage"
)

type Memory struct {
	data sync.Map
}

type Item struct {
	key      string
	value    string
	expireAt time.Time
	count    int
}

func (p *Memory) Clean(logic func(int64) bool) {
	var data []*Item
	p.data.Range(func(key, value interface{}) bool {
		data = append(data, value.(*Item))
		return true
	})

	pie.Any(pie.SortUsing(data, func(a, b *Item) bool {
		return a.count < b.count
	}), func(value *Item) bool {
		if !value.expireAt.IsZero() && value.expireAt.Before(time.Now()) {
			p.data.Delete(value)
			if logic(int64(len(value.value))) {
				return true
			}
		}
		return false
	})
}

func (p *Memory) SetNx(key string, value string) (bool, error) {
	_, loaded := p.data.LoadOrStore(key, &Item{
		key:   key,
		value: value,
	})
	return !loaded, nil
}

func (p *Memory) Get(key string) (string, error) {
	val, ok := p.data.Load(key)
	if !ok {
		return "", storage.ErrKeyNotExist
	}
	item := val.(*Item)
	if !item.expireAt.IsZero() && item.expireAt.Before(time.Now()) {
		p.data.Delete(key)
		return "", storage.ErrKeyNotExist
	}

	item.count++
	return item.value, nil
}

func (p *Memory) Expire(key string, timeout time.Duration) error {
	val, ok := p.data.Load(key)
	if !ok {
		return storage.ErrKeyNotExist
	}
	item := val.(*Item)
	if !item.expireAt.IsZero() && item.expireAt.Before(time.Now()) {
		p.data.Delete(key)
		return storage.ErrKeyNotExist
	}

	item.expireAt = time.Now().Add(timeout)
	p.data.Store(key, item)

	return nil
}

func (p *Memory) Close() error {
	return nil
}

func New() *Memory {
	p := &Memory{}
	return p
}
