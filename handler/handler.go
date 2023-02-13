package handler

import (
	"bufio"
	"bytes"
	"fmt"
	"time"

	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/valyala/fasthttp"
	"github.com/xihui-forever/goon/ctx"
)

type Handler struct {
	trie *Trie

	OnError func(ctx *ctx.Ctx, err error)
}

func NewHandler() *Handler {
	p := &Handler{
		trie: NewTrie(),
	}

	return p
}

func (p *Handler) CallOneOff(response *fasthttp.Response, request *fasthttp.Request) error {
	ctx, err := ctx.NewCtx(response, request)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return p.Call(ctx)
}

func (p *Handler) CallChrunked(response *fasthttp.Response, request *fasthttp.Request) error {
	ctx, err := ctx.NewCtx(response, request)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	ctx.Chucked(func(w *bufio.Writer) {
		p.Call(ctx)
	})
	return nil
}

func (p *Handler) Call(ctx *ctx.Ctx) error {
	defer func() {
		// 接收panic的信息，防止某一个请求导致程序崩溃
		if err := recover(); err != nil {
			log.Errorf("PANIC err:%v", err)
			if p.OnError != nil {
				p.OnError(ctx, fmt.Errorf("%v", err))
			}
		}
	}()

	var b bytes.Buffer
	b.WriteString("method:")
	b.WriteString(string(ctx.Method()))
	b.WriteString(" path:")
	b.WriteString(ctx.Path())

	log.Info("request ", b.String())
	defer func() {
		b.WriteString("used:")
		b.WriteString(time.Since(ctx.createdAt).String())
		log.Info("response ", b.String())
	}()

	itemList, err := p.trie.Find(ctx.Method(), ctx.Path())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	call := func(value *Item, ctx *ctx.Ctx) (res []byte, err error) {
		defer utils.CachePanicWithHandle(func(err interface{}) {
			if e, ok := err.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("%v", err)
			}
		})

		if ctx.isChuncked {
			err = value.TransferChrunked(ctx)
		} else {
			res, err = value.TransferOneOff(ctx)
		}
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	var key int
	var res []byte
	for index, value := range itemList {
		if value.method != PostUse {
			res, err = call(value, ctx)
			if err != nil {
				break
			}
		}
		key = index
	}

	for i := key; i >= 0; i-- {
		value := itemList[i]
		if value.method == PostUse {
			res, err = call(value, ctx)
			if err != nil {
				break
			}
		}
	}

	if !(ctx.isChuncked) {
		ctx.Write(res)
	}
	return nil
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

func (p *Handler) onError(ctx *ctx.Ctx, err error) {
	if p.OnError != nil {
		p.onError(ctx, err)
	}
}
