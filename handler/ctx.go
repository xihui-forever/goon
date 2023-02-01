package handler

import (
	"bufio"
	"time"

	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/valyala/fasthttp"
)

type Ctx struct {
	response   *fasthttp.Response
	request    *fasthttp.Request
	body       []byte
	method     Method
	path       string
	createdAt  time.Time
	isChuncked bool
}

func NewCtx(response *fasthttp.Response, request *fasthttp.Request) (*Ctx, error) {
	p := &Ctx{
		createdAt: time.Now().Truncate(time.Second),

		response: response,
		request:  request,

		body:       request.Body(),
		method:     Method(request.Header.Method()),
		path:       request.URI().String(),
		isChuncked: false,
	}

	return p, nil
}

func (p *Ctx) Method() Method {
	return p.method
}

func (p *Ctx) Path() string {
	return p.path
}

func (p *Ctx) Body() []byte {
	return p.body
}

func (p *Ctx) ParseBody(obj any) error {
	return sonic.Unmarshal(p.Body(), obj)
}

func (p *Ctx) GetSid() string {
	if p == nil {
		return ""
	}

	return string(p.request.Header.Peek("X-Session-Id"))
}

func (p *Ctx) Write(res []byte) {
	p.response.AppendBody(res)
}

func (p *Ctx) Send(res []byte) {
	// TODO: chunked
	p.response.AppendBody(res)
}

func (p *Ctx) Chucked(logic func(w *bufio.Writer)) {
	p.SetHeader(TransferEncoding, "chunked")
	p.response.SetBodyStreamWriter(func(w *bufio.Writer) {
		logic(w)
		err := w.Flush()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})
}

func (p *Ctx) SetHeader(key string, value string) {
	p.response.Header.Set(key, value)
	switch key {
	case TransferEncoding:
		switch value {
		case "chunked":
			// TODO ctx添加标记
			p.isChuncked = true
		}
	}
}

func (p *Ctx) Text(res string) {
	// TODO: 缓存数据
	p.response.Header.Set("Content-Type", "text/plain")
	p.response.AppendBodyString(res)
}

func (p *Ctx) Json(res any) error {
	// TODO: 缓存数据
	data, err := sonic.Marshal(res)
	if err != nil {
		return err
	}

	p.response.Header.Set("Content-Type", "application/json")
	p.response.AppendBody(data)

	return nil
}

func (p *Ctx) Jsonp(res any) error {
	// TODO: 缓存数据
	data, err := sonic.Marshal(res)
	if err != nil {
		return err
	}

	p.response.Header.Set("Content-Type", "application/javascript")
	p.response.AppendBody(data)

	return nil
}
