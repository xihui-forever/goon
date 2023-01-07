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
	trie *Trie
	lock sync.RWMutex
}

func NewHandler() *Handler {
	p := &Handler{
		trie: new(Trie),
	}

	return p
}

func NewAnalyze(logic interface{}) Analysis {
	s := reflect.TypeOf(logic)
	if s.Kind() != reflect.Func {
		panic("parameter is not func")
	}
	var analysis Analysis
	analysis.fun = logic

	if s.NumIn() != 2 {
		panic("num of func in is not 2")
	}
	x := s.In(0)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		panic("first in is must *Ctx")
	}
	if x.Name() != "Ctx" {
		panic("first in is must *Ctx")
	}
	analysis.in = append(analysis.in, x)

	x = s.In(1)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		panic("second in is must struct")
	}
	analysis.in = append(analysis.in, x)

	if s.NumOut() != 2 {
		panic("num of func in is not 2")
	}
	x = s.Out(0)
	for x.Kind() == reflect.Ptr {
		x = x.Elem()
	}
	if x.Kind() != reflect.Struct {
		panic("outFirst in is must struct")
	}
	analysis.out = append(analysis.out, x)

	x = s.Out(1)
	if x.Name() != "error" {
		panic("outSecond in is must error")
	}
	analysis.out = append(analysis.out, x)

	return analysis

}

func (p *Handler) Register(path string, logic interface{}) {
	// 校验logic是否是一个函数，并且函数的入参和出参是否规范
	// 同时记录path对应的logic
	analysis := NewAnalyze(logic)
	p.trie.Insert(path, &analysis)
}

func (p *Handler) Call(writer http.ResponseWriter, request *http.Request) error {
	// 根据request的path，找到对应的logic，并且调用
	node, err := p.trie.Find(request.URL.Path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var logic *Analysis
	for _, v := range node.childrenLogic {
		logic = v.logic
	}

	buf, err := io.ReadAll(request.Body)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	arg2 := reflect.New(logic.in[1])

	err = json.Unmarshal(buf, arg2.Interface())

	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	v := reflect.ValueOf(logic.fun)
	arg1 := &Ctx{
		writer:  writer,
		request: request,
	}
	call := v.Call([]reflect.Value{reflect.ValueOf(arg1), arg2})

	resp, err1 := json.Marshal(call[0].Interface())
	if err1 != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if call[1].Interface() != nil {
		return call[1].Interface().(error)
	}

	_, err = writer.Write(resp)
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
	}
}

func NewTrieNode(char string, analysis *Analysis) *trieNode {
	return &trieNode{
		char:          char,
		logic:         analysis,
		isEnding:      false,
		children:      make(map[rune]*trieNode),
		childrenLogic: make(map[*Analysis]*trieNode),
	}
}

func (p *Handler) NewTrie() *Trie {
	trieNode := NewTrieNode("/", nil)
	return &Trie{trieNode}
}

func (t *Trie) Insert(word string, analysis *Analysis) error {
	node := t.root
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			value = NewTrieNode(string(code), nil)
			node.children[code] = value
		}
		node = value
	}
	value, ok := node.childrenLogic[analysis]
	if !ok {
		if !ok && len(node.childrenLogic) == 0 {
			value = NewTrieNode("", analysis)
			node.childrenLogic[analysis] = value
		} else {
			errors.New("logic is wrong")
		}
	}
	node = value
	node.isEnding = true
	return nil
}

func (t *Trie) Find(word string) (*trieNode, error) {
	node := t.root
	for _, code := range word {
		value, ok := node.children[code]
		if !ok {
			return nil, errors.New("Path is not unRegistered")
		}
		node = value
	}
	return node, nil
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
		char          string
		logic         *Analysis
		isEnding      bool
		children      map[rune]*trieNode
		childrenLogic map[*Analysis]*trieNode
	}
	Trie struct {
		root *trieNode
	}
)

func main() {
	mux := NewHandler()
	mux.Register("/GetUser", func(ctx *Ctx, req *GetUserReq) (*GetUserRsp, error) {

		return &GetUserRsp{}, nil
	})

	mux.Register("/SetUser", func(ctx *Ctx, req *SetUserReq) (*SetUserRsp, error) {

		return nil, errors.New("not handle")
	})

	_ = http.ListenAndServe(":8080", mux)
}
