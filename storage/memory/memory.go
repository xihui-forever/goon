package memory

import (
	"sync"
	"time"

	"github.com/xihui-forever/goon/storage"
)

/*
		TODO: 1. 通过协程清理过期的key
	          2. 在用户调用方法时，清理过期的key
	          3. 允许用户选择清理的机制（1或者是2）
*/
type Memory struct {
	data sync.Map
}

type Item struct {
	value    string
	expireAt time.Time
}

func (m *Memory) SetNx(key string, value interface{}) (bool, error) {
	_, loaded := m.data.LoadOrStore(key, &Item{
		value: value.(string),
	})
	return !loaded, nil
}

func (m *Memory) Get(key string) (string, error) {
	val, ok := m.data.Load(key)
	if !ok {
		return "", storage.ErrKeyNotExist
	}
	item := val.(*Item)
	if !item.expireAt.IsZero() && item.expireAt.Before(time.Now()) {
		m.data.Delete(key)
		return "", storage.ErrKeyNotExist
	}

	return item.value, nil
}

func (m *Memory) Expire(key string, timeout time.Duration) error {
	val, ok := m.data.Load(key)
	if !ok {
		return storage.ErrKeyNotExist
	}
	item := val.(*Item)
	if !item.expireAt.IsZero() && item.expireAt.Before(time.Now()) {
		m.data.Delete(key)
		return storage.ErrKeyNotExist
	}

	item.expireAt = time.Now().Add(timeout)
	m.data.Store(key, item)

	return nil
}

func (m *Memory) Close() error {
	return nil
}

func New() *Memory {
	p := &Memory{}
	return p
}
