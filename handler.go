package goon

import (
	"bytes"
	"fmt"
	"time"

	"github.com/darabuchi/log"
	"github.com/elliotchance/pie/v2"
	"github.com/valyala/fasthttp"
)

type Handler func(ctx *Ctx) error

type App struct {
	trie *Trie

	OnError func(ctx *Ctx, err error)
}

func New() *App {
	p := &App{
		trie: NewTrie(),
	}

	return p
}

// func (p *Handler) CallChrunked(response *fasthttp.Response, request *fasthttp.Request) error {
// 	ctx, err := ctx.NewCtx(response, request)
// 	if err != nil {
// 		log.Errorf("err:%v", err)
// 		return err
// 	}
// 	ctx.Chucked(func(w *bufio.Writer) {
// 		p.Call(ctx)
// 	})
// 	return nil
// }

func (p *App) Call(context *fasthttp.RequestCtx) error {
	c := NewCtx(context)
	defer c.Close()

	defer func() {
		// 接收panic的信息，防止某一个请求导致程序崩溃
		if err := recover(); err != nil {
			log.Errorf("PANIC err:%v", err)
			if p.OnError != nil {
				p.OnError(c, fmt.Errorf("%v", err))
			}
		}
	}()

	var b bytes.Buffer
	b.WriteString("method:")
	b.WriteString(string(c.Method()))
	b.WriteString(" path:")
	b.WriteString(c.Path())

	log.Info("request ", b.String())
	defer func() {
		b.WriteString("used:")
		b.WriteString(time.Since(c.CreatedAt()).String())
		log.Info("response ", b.String())
	}()

	itemList, err := p.trie.Find(c.Method(), c.Path())
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	pie.Each(itemList, func(logic *Logic) {
		c.AppendHandler(logic.Handler)
	})

	return c.Next()
}

func (p *App) Register(method Method, path string, logic interface{}) {
	err := p.trie.Insert(path, NewLogic(method, logic))
	if err != nil {
		panic(fmt.Errorf("insert err:%v", err))
	}
}

func (p *App) Get(path string, logic any) {
	p.Register(Get, path, logic)
}

func (p *App) Post(path string, logic any) {
	p.Register(Post, path, logic)
}

func (p *App) Head(path string, logic any) {
	p.Register(Head, path, logic)
}

func (p *App) Use(path string, logic any) {
	p.Register(Use, path, logic)
}

func (p *App) PreUse(path string, logic any) {
	p.Register(PreUse, path, logic)
}

func (p *App) PostUse(path string, logic any) {
	p.Register(PostUse, path, logic)
}

func (p *App) onError(ctx *Ctx, err error) {
	if p.OnError != nil {
		p.onError(ctx, err)
	}
}
