package goon

func (p *Ctx) GetReqHeader(key string) string {
	if p == nil {
		return ""
	}

	return string(p.context.Request.Header.Peek(key))
}

func (p *Ctx) GetResHeader(key string) string {
	if p == nil {
		return ""
	}

	return string(p.context.Response.Header.Peek(key))
}

func (p *Ctx) SetReqHeader(key string, value string) {
	if p == nil {
		return
	}

	p.context.Request.Header.Set(key, value)
}

func (p *Ctx) SetResHeader(key string, value string) {
	if p == nil {
		return
	}

	p.context.Response.Header.Set(key, value)
}

func (p *Ctx) GetReqHeaderAll() map[string]string {
	if p == nil {
		return nil
	}

	m := map[string]string{}
	p.context.Request.Header.VisitAll(func(key, value []byte) {
		m[string(key)] = string(value)
	})

	return m
}

func (p *Ctx) GetResHeaderAll() map[string]string {
	if p == nil {
		return nil
	}

	m := map[string]string{}
	p.context.Response.Header.VisitAll(func(key, value []byte) {
		m[string(key)] = string(value)
	})

	return m
}
