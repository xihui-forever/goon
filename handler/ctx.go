package handler

import (
	"net/http"
)

type Ctx struct {
	writer  http.ResponseWriter
	request *http.Request
}

func (ctx *Ctx) GetSid() string {
	if ctx == nil {
		return ""
	}

	return ctx.request.Header.Get("X-Session-Id")
}
