package handler

import (
	"net/http"
)

type Ctx struct {
	writer  http.ResponseWriter
	request *http.Request
}

func (ctx *Ctx) Next() error {
	//是否需要拦截器
	//拦截器调用

	//if strings.HasPrefix(, "/User")
	return nil
}
