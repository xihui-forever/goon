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

	logicType := reflect.TypeOf(logic)
	if logicType.Kind() != reflect.Func {
		panic("parameter is not func")
	}

	if logicType.NumIn() == 0 && logicType.NumOut() == 0 {
		panic("num of func in is not normal ")
	}

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

	// 第二个参数必须是结构体
	if logicType.NumIn() == 2 {
		x = logicType.In(1)
		for x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() != reflect.Struct {
			panic("second in is must struct")
		}
		item.reqType = logicType.In(1)
		if logicType.NumOut() == 2 {
			x = logicType.Out(0).Elem()
			if x.Kind() == reflect.Ptr {
				x = x.Elem()
			}
			if x.Kind() != reflect.Struct {
				panic("outFirst in is must struct")
			}
			item.respType = logicType.Out(0)

			x = logicType.Out(1)
			if x.Name() != "error" {
				panic("outSecond in is must error")
			}
		}
	}

	if logicType.NumOut() == 2 {
		x = logicType.Out(0)
		if x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() == reflect.Ptr {
			x = x.Elem()
		}
		if x.Kind() != reflect.Struct {
			panic("outFirst in is must struct")
		}
		item.respType = logicType.Out(0)

		x = logicType.Out(1)
		if x.Name() != "error" {
			panic("outSecond in is must error")
		}
		return item

	}

	x = logicType.Out(0)
	if x.Name() != "error" {
		panic("out in is must error")
	}

	return item
}

func (p *Item) Call(writer http.ResponseWriter, request *http.Request) error {
	buf, err := io.ReadAll(request.Body)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var res []reflect.Value
	if p.reqType == nil {
		res = p.logic.Call([]reflect.Value{
			reflect.ValueOf(&Ctx{
				writer:  writer,
				request: request,
			}),
		})
	} else {
		req := reflect.New(p.reqType)
		err = json.Unmarshal(buf, req.Interface())
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		res = p.logic.Call([]reflect.Value{
			reflect.ValueOf(&Ctx{
				writer:  writer,
				request: request,
			}),
			req,
		})
	}
	if len(res) == 1 {
		if res[0].Interface() != nil {
			return res[1].Interface().(error)
		}
		return nil
	}
	if res[1].Interface() != nil {
		return res[1].Interface().(error)
	}

	var resp []byte
	if res[0].IsValid() {
		resp, err = json.Marshal(res[0].Interface())
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
