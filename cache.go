package goon

import (
	"sync"

	"github.com/darabuchi/utils"
)

// Cache 缓存的数据，用于对 `Ctx` 的缓存
type Cache struct {
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

func (p *Cache) GetUint64(key string) (uint64, bool) {
	v, ok := p.Get(key)
	if !ok {
		return 0, false
	}

	return utils.ToUint64(v), true
}

func (p *Cache) GetUint32(key string) (uint32, bool) {
	v, ok := p.Get(key)
	if !ok {
		return 0, false
	}

	return utils.ToUint32(v), true
}

func (p *Cache) GetInt64(key string) (int64, bool) {
	v, ok := p.Get(key)
	if !ok {
		return 0, false
	}

	return utils.ToInt64(v), true
}

func (p *Cache) GetInt32(key string) (int32, bool) {
	v, ok := p.Get(key)
	if !ok {
		return 0, false
	}

	return utils.ToInt32(v), true
}

func (p *Cache) GetInt(key string) (int, bool) {
	v, ok := p.Get(key)
	if !ok {
		return 0, false
	}

	return utils.ToInt(v), true
}

func (p *Cache) GetFloat64(key string) (float64, bool) {
	v, ok := p.Get(key)
	if !ok {
		return 0, false
	}

	return utils.ToFloat64(v), true
}

func (p *Cache) GetFloat32(key string) (float32, bool) {
	v, ok := p.Get(key)
	if !ok {
		return 0, false
	}

	return utils.ToFloat32(v), true
}

func (p *Cache) GetBool(key string) (bool, bool) {
	v, ok := p.Get(key)
	if !ok {
		return false, false
	}

	return utils.ToBool(v), true
}

func (p *Cache) GetString(key string) (string, bool) {
	v, ok := p.Get(key)
	if !ok {
		return "", false
	}

	return utils.ToString(v), true
}

func (p *Cache) GetUint64Def(key string, def uint64) uint64 {
	v, ok := p.GetUint64(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetUint32Def(key string, def uint32) uint32 {
	v, ok := p.GetUint32(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetInt64Def(key string, def int64) int64 {
	v, ok := p.GetInt64(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetInt32Def(key string, def int32) int32 {
	v, ok := p.GetInt32(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetIntDef(key string, def int) int {
	v, ok := p.GetInt(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetFloat64Def(key string, def float64) float64 {
	v, ok := p.GetFloat64(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetFloat32Def(key string, def float32) float32 {
	v, ok := p.GetFloat32(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetBoolDef(key string, def bool) bool {
	v, ok := p.GetBool(key)
	if !ok {
		return def
	}

	return v
}

func (p *Cache) GetStringDef(key string, def string) string {
	v, ok := p.GetString(key)
	if !ok {
		return def
	}

	return v
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
