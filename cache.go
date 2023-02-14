package goon

import (
	"sync"
)

type Cache struct {
	// 缓存的数据

	m sync.Map
}

var cachePool = sync.Pool{
	New: func() interface{} {
		return newCache()
	},
}

func newCache() *Cache {
	p := &Cache{}

	return p
}

func NewCache() *Cache {
	return cachePool.Get().(*Cache)
}

func (p *Cache) Get(key string) (interface{}, bool) {
	return p.m.Load(key)
}

func (p *Cache) Set(key string, value interface{}) {
	p.m.Store(key, value)
}

func (p *Cache) Delete(key string) {
	p.m.Delete(key)
}

func (p *Cache) Reset() {
	p.m.Range(func(key, value interface{}) bool {
		p.m.Delete(key)
		return true
	})
}

func (p *Cache) Close() {
	p.Reset()
	cachePool.Put(p)
}
