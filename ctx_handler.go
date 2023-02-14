package goon

func (p *Ctx) AppendHandler(handler func(ctx *Ctx) error) {
	p.handlers = append(p.handlers, handler)
}

func (p *Ctx) Next() error {
	if p.handlerIdx >= len(p.handlers) {
		return nil
	}

	handler := p.handlers[p.handlerIdx]
	p.handlerIdx++

	return handler(p)
}
