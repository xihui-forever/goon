package handler

import (
	"net/http"
)

type Ctx struct {
	writer  http.ResponseWriter
	request *http.Request
}

func (p *Ctx) GetSid() string {
	if p == nil {
		return ""
	}

	return p.request.Header.Get("X-Session-Id")
}
