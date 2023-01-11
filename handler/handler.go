package handler

import (
	"fmt"
	"github.com/darabuchi/log"
	"net/http"
	"strings"
)

type Handler struct {
	trie *Trie
}

func NewHandler() *Handler {
	p := &Handler{
		trie: NewTrie(),
	}

	return p
}

func (p *Handler) IsValid(writer http.ResponseWriter, request *http.Request) error {
	if strings.HasPrefix(request.URL.Path, "/User/") {
		request.Method = ""
		request.URL.Path = "/User"

		err := p.Call(writer, request)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
	}
	return nil
}
func (p *Handler) Call(writer http.ResponseWriter, request *http.Request) error {
	// 根据request的path，找到对应的logic，并且调用
	method := Method(request.Method)

	item, err := p.trie.Find(method, request.URL.Path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	log.Info(item.reqType)
	err = item.Call(writer, request)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (p *Handler) Register(method Method, path string, logic interface{}) {
	err := p.trie.Insert(method, path, NewItem(logic))
	if err != nil {
		panic(fmt.Errorf("insert err:%v", err))
	}
}

func (p *Handler) Get(path string, logic any) {
	p.Register(Get, path, logic)
}

func (p *Handler) Post(path string, logic any) {
	p.Register(Post, path, logic)
}

func (p *Handler) Head(path string, logic any) {
	p.Register(Head, path, logic)
}

func (p *Handler) Use(path string, logic any) {
	p.Register("", path, logic)
}

func (p *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	method := request.Method
	path := request.URL.Path
	if p.IsValid(writer, request) == nil {
		request.Method = method
		request.URL.Path = path
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
}
