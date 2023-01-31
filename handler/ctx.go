package handler

import (
	"github.com/valyala/fasthttp"
)

type Ctx struct {
	response *fasthttp.Response
	request  *fasthttp.Request
}

func (p *Ctx) GetSid() string {
	if p == nil {
		return ""
	}

	return string(p.request.Header.Peek("X-Session-Id"))
}

func (p *Ctx) Write(res []byte) error {
	p.response.AppendBody(res)
	return nil
}
