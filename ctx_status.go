package goon

func (p *Ctx) SetStatusCode(code int) {
	p.context.Response.SetStatusCode(code)
}
