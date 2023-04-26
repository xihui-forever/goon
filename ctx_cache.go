package goon

func (p *Ctx) Cache() *Cache {
	return p.cache
}

func (p *Ctx) Get(key string) (interface{}, bool) {
	return p.cache.Get(key)
}

func (p *Ctx) GetWithDef(key string, def interface{}) interface{} {
	val, ok := p.Get(key)
	if ok {
		return val
	}
	return def
}

func (p *Ctx) Set(key string, value interface{}) {
	p.cache.Set(key, value)
}
