package lifo

import (
	"sync"
	"time"

	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon/storage"
)

type Memory struct {
	data []*Item

	lock sync.RWMutex
}

type Item struct {
	key      string
	value    string
	expireAt time.Time
}

func (p *Memory) Clean(logic func(size int64) bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	for k, v := range p.data {
		if !v.expireAt.IsZero() && v.expireAt.Before(time.Now()) {
			p.data = append(p.data[k:], p.data[k+1:]...)
			if logic(int64(len(v.value))) {
				return
			}
		}
	}
}

func (p *Memory) visit(key string) (*Item, error) {
	for k, v := range p.data {
		if key == v.key {
			if !v.expireAt.IsZero() && v.expireAt.Before(time.Now()) {
				p.data = append(p.data[k:], p.data[k+1:]...)
				return nil, storage.ErrKeyNotExist
			}
			return v, nil
		}
	}
	return nil, storage.ErrKeyNotExist
}

func (p *Memory) SetNx(key string, value string) (bool, error) {
	item := &Item{
		key:   key,
		value: value,
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	_, err := p.visit(key)
	if err != nil {
		if err == storage.ErrKeyNotExist {
			p.data = append(p.data, &Item{})
			copy(p.data[1:], p.data[0:])
			p.data[0] = item
			return true, nil
		}

		log.Errorf("err:%v", err)
		return false, err
	}

	return false, nil
}

func (p *Memory) Get(key string) (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	item, err := p.visit(key)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	return item.value, nil
}

func (p *Memory) Expire(key string, timeout time.Duration) error {
	item, err := p.visit(key)
	if err != nil {
		return err
	}

	item.expireAt = time.Now().Add(timeout)
	return nil
}

func (p *Memory) Close() error {
	return nil
}

func New() *Memory {
	p := &Memory{
		data: make([]*Item, 0, 1024),
	}

	return p
}
