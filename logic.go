package goon

import (
	"encoding/json"
	"reflect"

	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
)

type Logic struct {
	logic    reflect.Value
	reqType  reflect.Type
	respType reflect.Type
	method   Method
}

func NewLogic(method Method, logic any) *Logic {
	item := &Logic{
		logic:  reflect.ValueOf(logic),
		method: method,
	}

	isStruct := func(x reflect.Type) bool {
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}

		return x.Kind() == reflect.Struct
	}

	isError := func(x reflect.Type) bool {
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}

		return x.Name() == "error"
	}

	logicType := reflect.TypeOf(logic)
	if logicType.Kind() != reflect.Func {
		panic("parameter is not func")
	}

	// 不存在第一个入参
	if logicType.NumIn() == 0 {
		panic("num of func in can not be empty")
	}

	// 存在第一个入参
	// 第一个参数必须是*Ctx
	x := logicType.In(0)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}

	if x.Kind() != reflect.Struct {
		panic("first in is must *Ctx")
	}
	if x.Name() != "Ctx" {
		panic("first in is must *Ctx")
	}

	// 存在第二个入参
	// 第二个参数必须是结构体
	if logicType.NumIn() > 1 {
		x = logicType.In(1)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		item.reqType = x
		if !isStruct(item.reqType) {
			panic("2rd in is must struct")
		}
	}

	switch logicType.NumOut() {
	case 1:
		// 只有一个出参数
		if !isError(logicType.Out(0)) {
			panic("1st out must error")
		}
	case 2:
		// 有两个返回
		item.respType = logicType.Out(0)
		if !isStruct(item.respType) {
			panic("1st out must struct")
		}

		if !isError(logicType.Out(1)) {
			panic("2rd out must error")
		}
	case 0:
		fallthrough
	default:
		panic("num of func out must 1 or 2")
	}

	return item
}

func (p *Logic) Method() Method {
	return p.method
}

func (p *Logic) TransferChrunked(ctx *Ctx) error {
	in := []reflect.Value{
		// 第一个入参是固定的
		reflect.ValueOf(ctx),
	}

	// 如果存在第二个入参
	if p.reqType != nil {
		req := reflect.New(p.reqType)
		err := ctx.ParseBody(req.Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		in = append(in, req)
	}

	// 调用处理方法
	out := p.logic.Call(in)

	// 只有一个返回值的
	if p.respType == nil {
		if out[0].Interface() != nil {
			return out[0].Interface().(error)
		}
		return nil
	}

	// 有两个返回值的
	if out[1].Interface() != nil {
		return out[1].Interface().(error)
	}

	var resp []byte
	var err error
	if out[0].IsValid() {
		resp, err = json.Marshal(out[0].Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	} else {
		resp = []byte("{}")
	}

	ctx.Send(string(resp))

	return nil
}

func (p *Logic) Transfer(ctx *Ctx) ([]byte, error) {
	in := []reflect.Value{
		// 第一个入参是固定的
		reflect.ValueOf(ctx),
	}

	// 如果存在第二个入参
	if p.reqType != nil {
		req := reflect.New(p.reqType)
		err := ctx.ParseBody(req.Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		in = append(in, req)
	}

	// 调用处理方法
	out := p.logic.Call(in)

	// 只有一个返回值的
	if p.respType == nil {
		if out[0].Interface() != nil {
			return nil, out[0].Interface().(error)
		}
		return nil, nil
	}

	// 有两个返回值的
	if out[1].Interface() != nil {
		return nil, out[1].Interface().(error)
	}

	var resp string
	var err error
	if out[0].IsValid() {
		resp, err = sonic.MarshalString(out[0].Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}
	} else {
		resp = ("{}")
	}

	ctx.Send(resp)

	return []byte(resp), nil
}

func (p *Logic) Handler(ctx *Ctx) error {
	in := []reflect.Value{
		// 第一个入参是固定的
		reflect.ValueOf(ctx),
	}

	// 如果存在第二个入参
	if p.reqType != nil {
		req := reflect.New(p.reqType)
		err := ctx.ParseBody(req.Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		in = append(in, req)
	}

	// 调用处理方法
	out := p.logic.Call(in)

	// 只有一个返回值的
	if p.respType == nil {
		if out[0].Interface() != nil {
			return out[0].Interface().(error)
		}
		return nil
	}

	// 有两个返回值的
	if out[1].Interface() != nil {
		return out[1].Interface().(error)
	}

	if out[0].IsValid() {
		if ctx.o.PermHeader == "" {
			return ctx.Json(out[0].Interface())
		}

		return ctx.JsonWithPerm(ctx.GetReqHeader(ctx.o.PermHeader), out[0].Interface())
	}

	return ctx.Send("{}")
}
