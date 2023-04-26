package goon

import (
	"fmt"
	"github.com/darabuchi/log"
	"github.com/elliotchance/pie/v2"
	"github.com/valyala/fasthttp"
	"runtime/debug"
)

type Handler func(ctx *Ctx) error

type App struct {
	trie *Trie

	OnError func(ctx *Ctx, err error)

	o *option
}

var a = New()

func New(opts ...*option) *App {
	p := &App{
		trie: NewTrie(),
		o:    &option{},
	}

	p.WithOptions(opts...)

	return p
}

func (p *App) WithOptions(opts ...*option) *App {
	p.o = p.o.applyOption(opts...)
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
	c := NewCtx(context, p.o)
	defer c.Close()

	defer func() {
		// 接收panic的信息，防止某一个请求导致程序崩溃
		if err := recover(); err != nil {
			log.Errorf("PANIC err:%v", err)
			st := debug.Stack()
			if len(st) > 0 {
				log.Error(string(st))
			} else {
				log.Errorf("stack is empty (%s)", err)
			}
			if p.OnError != nil {
				p.OnError(c, fmt.Errorf("%v", err))
			}
		}
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
	log.Infof("register %s %s", method, path)
	err := p.trie.Insert(path, NewLogic(method, logic))
	if err != nil {
		panic(fmt.Errorf("insert err:%v", err))
	}
}

func (p *App) Get(path string, logic any) {
	p.Register(MethodGet, path, logic)
}

func (p *App) Post(path string, logic any) {
	p.Register(MethodPost, path, logic)
}

func (p *App) Head(path string, logic any) {
	p.Register(MethodHead, path, logic)
}

func (p *App) Put(path string, logic any) {
	p.Register(MethodPut, path, logic)
}

func (p *App) Option(path string, logic any) {
	p.Register(MethodOption, path, logic)
}

func (p *App) Delete(path string, logic any) {
	p.Register(MethodDelete, path, logic)
}

func (p *App) Use(path string, logic any) {
	p.Register(MethodUse, path, logic)
}

func (p *App) PreUse(path string, logic any) {
	p.Register(MethodPreUse, path, logic)
}

func (p *App) PostUse(path string, logic any) {
	p.Register(MethodPostUse, path, logic)
}

func (p *App) onError(ctx *Ctx, err error) {
	if p.OnError != nil {
		p.onError(ctx, err)
	}
}

func (p *App) ListenAndServe(addr string) error {
	return fasthttp.ListenAndServe(addr, func(ctx *fasthttp.RequestCtx) {
		log.SetTrace(log.GenTraceId())
		defer log.DelTrace()

		ctx.Response.Header.Set("X-Trace-Id", log.GetTrace())

		ctx.Response.Header.SetStatusCode(fasthttp.StatusOK)
		err := p.Call(ctx)
		if err != nil {
			log.Errorf("err:%v", err)
			_, e := ctx.Write([]byte(err.Error()))
			if e != nil {
				log.Errorf("err:%v", e)
			}
		}
	})
}

func Register(method Method, path string, logic interface{}) {
	a.Register(method, path, logic)
}

func Get(path string, logic any) {
	a.Get(path, logic)
}

func Post(path string, logic any) {
	a.Post(path, logic)
}

func Head(path string, logic any) {
	a.Head(path, logic)
}

func Put(path string, logic any) {
	a.Put(path, logic)
}

func Option(path string, logic any) {
	a.Option(path, logic)
}

func Delete(path string, logic any) {
	a.Delete(path, logic)
}

func Use(path string, logic any) {
	a.Use(path, logic)
}

func PreUse(path string, logic any) {
	a.PreUse(path, logic)
}

func PostUse(path string, logic any) {
	a.PostUse(path, logic)
}

func ListenAndServe(addr string) error {
	return a.ListenAndServe(addr)
}
