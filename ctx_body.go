package goon

import (
	"reflect"
	"strings"

	"github.com/bytedance/sonic"
	"github.com/elliotchance/pie/v2"
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

func (p *Ctx) Send(res string) error {
	p.respBody = &res
	return nil
}

// func (p *Ctx) Chucked(logic func(w *bufio.Writer)) {
// 	p.SetResHeader(TransferEncoding, "chunked")
// 	p.context.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
// 		logic(w)
// 		err := w.Flush()
// 		if err != nil {
// 			log.Errorf("err:%v", err)
// 			return
// 		}
// 	})
// }

func (p *Ctx) Text(res string) error {
	// TODO: 缓存数据
	p.context.Response.Header.Set("Content-Type", "text/plain")
	return p.Send(res)
}

func (p *Ctx) Json(res any) error {
	// TODO: 缓存数据
	data, err := sonic.MarshalString(res)
	if err != nil {
		return err
	}

	p.context.Response.Header.Set("Content-Type", "application/json")

	return p.Send(data)
}

func (p *Ctx) JsonWithPerm(permStr string, res any) error {
	perms := strings.Split(permStr, "|")

	var logic func(res any) reflect.Value
	logic = func(res any) reflect.Value {
		v := reflect.ValueOf(res)
		t := reflect.TypeOf(res)

		for i := 0; i < v.NumField(); i++ {
			iv := v.Field(i)
			if !iv.IsValid() {
				continue
			}

			it := t.Field(i)

			tagPerm := it.Tag.Get("perm")
			if tagPerm == "" || tagPerm == "-" {
				switch it.Type.Kind() {
				case reflect.Ptr, reflect.Struct:
					iv.Set(logic(iv.Interface()))
				}

				continue
			}

			tagPerms := strings.Split(tagPerm, "|")

			if pie.Any(tagPerms, func(perm string) bool {
				return pie.Contains(perms, perm)
			}) {

				switch it.Type.Kind() {
				case reflect.Ptr, reflect.Struct:
					iv.Set(logic(iv.Interface()))
				}

				continue
			}

			iv.Set(reflect.Zero(iv.Type()))
		}

		return v
	}

	// TODO: 缓存数据
	data, err := sonic.MarshalString(logic(res).Interface())
	if err != nil {
		return err
	}

	p.context.Response.Header.Set("Content-Type", "application/json")

	return p.Send(data)
}

func (p *Ctx) Jsonp(res any) error {
	// TODO: 缓存数据
	data, err := sonic.MarshalString(res)
	if err != nil {
		return err
	}

	p.context.Response.Header.Set("Content-Type", "application/javascript")

	return p.Send(data)
}
