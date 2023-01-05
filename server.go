package main

import (
	"encoding/json"
	"errors"
	"github.com/darabuchi/log"
	"io"
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
	lock   sync.RWMutex
}

func NewHandler() *Handler {
	p := &Handler{
		mapper: map[string]interface{}{},
	}

	return p
}

func (p *Handler) MapPathToFunc(path string, s interface{}) {
	p.lock.Lock()
	if _, ok := p.mapper[path]; !ok {
		p.mapper[path] = s
	}
	p.lock.Unlock()
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
		panic("second in is must struct")
	}

	if s.NumOut() != 2 {
		panic("num of func in is not 2")
	}
	x = s.Out(0).Elem()
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		panic("outFirst in is must struct")
	}

	x = s.Out(1)
	if x.Name() != "error" {
		panic("outSecond in is must error")
	}

	p.MapPathToFunc(path, s)
}

func (p *Handler) Call(writer http.ResponseWriter, request *http.Request) error {
	// 根据request的path，找到对应的logic，并且调用
	logic, ok := p.mapper[request.URL.Path]
	if !ok {
		return errors.New("Path is unregistered")
	}
	s := reflect.TypeOf(logic)
	in1 := s.In(0).Elem()
	in2 := s.In(1).Elem()
	log.Info(in2)

	buf, err := io.ReadAll(request.Body)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = json.Unmarshal(buf, in2)
	log.Info(in2)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	call := reflect.New(s).Call([]reflect.Value{reflect.ValueOf(in1), reflect.ValueOf(in2)})
	log.Info(call)
	resp, err := json.Marshal(call[0].Interface())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_, err = writer.Write(resp)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
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
	mux.Register("/GetUser", func(ctx *Ctx, req *GetUserReq) (*GetUserRsp, error) {

		return nil, errors.New("not handle")
	})

	mux.Register("/SetUser", func(ctx *Ctx, req *SetUserReq) (*SetUserRsp, error) {

		return nil, errors.New("not handle")
	})
}
