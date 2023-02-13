package ctx

func (p *Ctx) GetReqHeader(key string) string {
	if p == nil {
		return ""
	}

	return string(p.request.Header.Peek(key))
}

func (p *Ctx) GetResHeader(key string) string {
	if p == nil {
		return ""
	}

	return string(p.response.Header.Peek(key))
}

func (p *Ctx) SetReqHeader(key string, value string) {
	if p == nil {
		return
	}

	p.request.Header.Set(key, value)
}

func (p *Ctx) SetResHeader(key string, value string) {
	if p == nil {
		return
	}

	p.response.Header.Set(key, value)
}

func (p *Ctx) GetReqHeaderAll() map[string]string {
	if p == nil {
		return nil
	}

	m := map[string]string{}
	p.request.Header.VisitAll(func(key, value []byte) {
		m[string(key)] = string(value)
	})

	return m
}

func (p *Ctx) GetResHeaderAll() map[string]string {
	if p == nil {
		return nil
	}

	m := map[string]string{}
	p.response.Header.VisitAll(func(key, value []byte) {
		m[string(key)] = string(value)
	})

	return m
}

func (p *Ctx) GetSid() string {
	if p == nil {
		return ""
	}

	return string(p.request.Header.Peek("X-Session-Id"))
}
