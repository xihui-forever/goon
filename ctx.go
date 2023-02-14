package goon

import (
	"bufio"
	"time"

	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/valyala/fasthttp"
)

type Ctx struct {
	context *fasthttp.RequestCtx

	createdAt time.Time

	handlerIdx int
	handlers   []func(ctx *Ctx) error
	cache      *Cache
}

func NewCtx(context *fasthttp.RequestCtx) *Ctx {
	p := &Ctx{
		createdAt: time.Now().Truncate(time.Second),
		context:   context,
		cache:     NewCache(),
	}

	return p
}

func (p *Ctx) Context() *fasthttp.RequestCtx {
	return p.context
}

func (p *Ctx) Method() Method {
	return Method(p.context.Request.Header.Method())
}

func (p *Ctx) Path() string {
	return p.context.Request.URI().String()
}

func (p *Ctx) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Ctx) Close() {
	p.cache.Close()
}

func (p *Ctx) Body() []byte {
	return p.context.Request.Body()
}

func (p *Ctx) ParseBody(obj any) error {
	return sonic.Unmarshal(p.Body(), obj)
}

func (p *Ctx) Write(res []byte) {
	p.context.Response.AppendBody(res)
}

func (p *Ctx) Send(res []byte) error {
	p.context.Response.AppendBody(res)
	return nil
}

func (p *Ctx) Chucked(logic func(w *bufio.Writer)) {
	p.SetHeader(TransferEncoding, "chunked")
	p.context.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
		logic(w)
		err := w.Flush()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})
}

func (p *Ctx) SetHeader(key string, value string) {
	p.context.Response.Header.Set(key, value)
	switch key {
	case TransferEncoding:
		switch value {
		case "chunked":
			// TODO ctx添加标记
		}
	}
}

func (p *Ctx) Text(res string) {
	// TODO: 缓存数据
	p.context.Response.Header.Set("Content-Type", "text/plain")
	p.context.Response.AppendBodyString(res)
}

func (p *Ctx) Json(res any) error {
	// TODO: 缓存数据
	data, err := sonic.Marshal(res)
	if err != nil {
		return err
	}

	p.context.Response.Header.Set("Content-Type", "application/json")
	p.context.Response.AppendBody(data)

	return nil
}

func (p *Ctx) Jsonp(res any) error {
	// TODO: 缓存数据
	data, err := sonic.Marshal(res)
	if err != nil {
		return err
	}

	p.context.Response.Header.Set("Content-Type", "application/javascript")
	p.context.Response.AppendBody(data)

	return nil
}
