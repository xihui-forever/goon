package main

import (
	"errors"
	"net/http"
	"reflect"
	"sync"
)

type Ctx struct {
	writer  http.ResponseWriter
	request *http.Request
}

type Handler struct {
	mapper map[string]interface{}
}

func NewHandler() *Handler {
	p := &Handler{}

	return p
}

var lock sync.RWMutex

func (p *Handler) MapPathToFunc(path string, s interface{}) {
	lock.Lock()
	if _, ok := p.mapper[path]; !ok {
		p.mapper[path] = s
	}
	lock.Unlock()
}

func (p *Handler) Register(path string, logic interface{}) {
	// 校验logic是否是一个函数，并且函数的入参和出参是否否则规范
	// 同时记录path对应的logic

	s := reflect.TypeOf(logic)
	if s.Kind() != reflect.Func {
		panic("parameter is not func")
	}
	if s.NumIn() != 2 {
		panic("num of func in is not 2")
	}
	for i := 0; i < s.NumIn(); i++ {
		if s.In(i) == nil {
			panic("null pointer exception")
		}
	}
	x := s.In(0).Elem()
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Kind() != reflect.Struct {
		panic("first in is must *Ctx")
	}

	if x.Name() != "Ctx" {
		panic("first in is must *Ctx")
	}

	x = s.In(1).Elem()
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		panic("second in is must *GetUserReq")
	}

	x = s.Out(1).Elem()
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		panic("outFirst in is must *GetUserReq")
	}

	x = s.Out(2)
	if x.Name() != "error" {
		panic("outSecond in is must error")
	}

	p.MapPathToFunc(path, s)

}

type (
	GetUserReq struct {
	}

	GetUserRsp struct {
	}
)

type (
	SetUserReq struct {
	}

	SetUserRsp struct {
	}
)

func main() {
	mux := NewHandler()
	mux.mapper = make(map[string]interface{})
	mux.Register("/GetUser", func(ctx *Ctx, req *GetUserReq) (*GetUserRsp, error) {

		return nil, errors.New("not handle")
	})

	mux.Register("/SetUser", func(ctx *Ctx, req *SetUserReq) (*SetUserRsp, error) {

		return nil, errors.New("not handle")
	})
}
