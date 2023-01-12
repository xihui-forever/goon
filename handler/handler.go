package handler

import (
	"fmt"
	"github.com/darabuchi/log"
	"net/http"
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

func (p *Handler) Call(writer http.ResponseWriter, request *http.Request) error {
	// 根据request的path，找到对应的logic，并且调用
	method := Method(request.Method)

	itemList, err := p.trie.Find(method, request.URL.Path)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	var key int
	for index, value := range itemList {
		if value.method == PreUse {
			err = value.Call(writer, request)
			if err != nil {
				break
			}
		}
		key = index
	}
	if key == len(itemList)-1 {
		err = itemList[len(itemList)-1].Call(writer, request)
	}

	for i := key; i >= 0; i-- {
		value := itemList[i]
		if value.method == PostUse {
			value.Call(writer, request)
		}
	}

	return err
}

func (p *Handler) Register(method Method, path string, logic interface{}) {
	err := p.trie.Insert(path, NewItem(method, logic))
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

func (p *Handler) PreUse(path string, logic any) {
	p.Register(PreUse, path, logic)
}

func (p *Handler) PostUse(path string, logic any) {
	p.Register(PostUse, path, logic)
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
