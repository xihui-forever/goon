package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

type Item struct {
	logic    reflect.Value
	reqType  reflect.Type
	respType reflect.Type
}

func NewItem(logic any) *Item {
	item := &Item{
		logic: reflect.ValueOf(logic),
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

		return x.Name() != "error"
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
		item.reqType = logicType.In(1)
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

func (p *Item) Call(writer http.ResponseWriter, request *http.Request) error {
	in := []reflect.Value{
		// 第一个入参是固定的
		reflect.ValueOf(&Ctx{
			writer:  writer,
			request: request,
		}),
	}

	// 如果存在第二个入参
	if p.reqType != nil {
		buf, err := io.ReadAll(request.Body)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		req := reflect.New(p.reqType)
		err = json.Unmarshal(buf, req.Interface())
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

	_, err = writer.Write(resp)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

type Handler struct {
	trie *Trie
	lock sync.RWMutex
}

func NewHandler() *Handler {
	p := &Handler{
		trie: new(Trie),
	}

	return p
}

func (p *Handler) Register(path string, logic interface{}) {
	// 校验logic是否是一个函数，并且函数的入参和出参是否规范
	// 同时记录path对应的logic

	err := p.trie.Insert(path, NewItem(logic))
	if err != nil {
		panic(fmt.Errorf("insert err:%v", err))
	}
}

func (p *Handler) Call(writer http.ResponseWriter, request *http.Request) error {
	// 根据request的path，找到对应的logic，并且调用
	item, err := p.trie.Find(request.URL.Path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = item.Call(writer, request)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (p *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := p.Call(writer, request)
	if err != nil {
		log.Errorf("err:%v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		_, e := writer.Write([]byte(err.Error()))
		if e != nil {
			log.Errorf("err:%v", e)
		}
	}
}

func NewTrieNode(char string, item *Item) *trieNode {
	return &trieNode{
		char:     char,
		logic:    item,
		isEnding: false,
		children: make(map[rune]*trieNode),
	}
}

func (p *Handler) NewTrie() *Trie {
	trieNode := NewTrieNode("/", nil)
	return &Trie{trieNode}
}

func (t *Trie) Insert(word string, item *Item) error {
	node := t.root
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			value = NewTrieNode(string(code), nil)
			node.children[code] = value
		}
		node = value
	}

	if node.logic != nil {
		return errors.New("logic already exist")
	}

	node.logic = item
	node.isEnding = true
	return nil
}

func (t *Trie) Find(word string) (*Item, error) {
	node := t.root
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			return nil, errors.New("path is not unRegistered")
		}
		node = value
	}
	return node.logic, nil
}

type Analysis struct {
	intNum int
	outNum int
	in     []reflect.Type
	out    []reflect.Type
	fun    interface{}
}

type (
	GetUserReq struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	GetUserRsp struct {
		A string `json:"a"`
		B string `json:"b"`
	}
)

type (
	SetUserReq struct {
		C string `json:"c"`
		D string `json:"d"`
	}

	SetUserRsp struct {
		C string `json:"c"`
		D string `json:"d"`
	}
)

type (
	trieNode struct {
		char     string
		logic    *Item
		isEnding bool
		children map[rune]*trieNode
	}

	Trie struct {
		root *trieNode
	}
)

func main() {
	mux := NewHandler()

	mux.Register("/Check", func(ctx *Ctx) error {
		return nil
	})

	mux.Register("/GetMe", func(ctx *Ctx) (**GetUserRsp, error) {

		return nil, errors.New("not handle")
	})

	mux.Register("/SetMe", func(ctx *Ctx, req *SetUserReq) error {

		return nil
	})

	mux.Register("/SetUser", func(ctx *Ctx, req *SetUserReq) (*SetUserRsp, error) {

		return nil, errors.New("not handle")
	})

	_ = http.ListenAndServe(":8080", mux)
}
