package goon

import (
	"bufio"

	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
)

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
	p.SetResHeader(TransferEncoding, "chunked")
	p.context.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
		logic(w)
		err := w.Flush()
		if err != nil {
			log.Errorf("err:%v", err)
			return
		}
	})
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
