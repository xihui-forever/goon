package handler

import (
	"fmt"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/valyala/fasthttp"
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

func (p *Handler) Call(response *fasthttp.Response, request *fasthttp.Request) error {
	// 根据request的path，找到对应的logic，并且调用
	method := Method(request.Header.Method())
	itemList, err := p.trie.Find(method, request.URI().String())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	context := &Ctx{
		response: response,
		request:  request,
	}

	call := func(value *Item, ctx *Ctx) (res []byte, err error) {
		defer utils.CachePanicWithHandle(func(err interface{}) {
			if e, ok := err.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", err)
			}
		})

		res, err = value.CallOne(ctx)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	var key int
	var res []byte
	for index, value := range itemList {
		if value.method != PostUse {
			res, err = call(value, context)
			if err != nil {
				break
			}
		}
		key = index
	}

	for i := key; i >= 0; i-- {
		value := itemList[i]
		if value.method == PostUse {
			res, err = call(value, context)
			if err != nil {
				break
			}
		}
	}
	return context.Write(res)

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
