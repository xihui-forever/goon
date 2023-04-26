package goon

import (
	"bytes"
	"github.com/darabuchi/log"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/valyala/fasthttp"
)

type Ctx struct {
	context *fasthttp.RequestCtx

	createdAt time.Time

	respBody *string // 返回的body
	o        *option

	handlerIdx int
	handlers   []func(ctx *Ctx) error
	cache      *Cache
}

func NewCtx(context *fasthttp.RequestCtx, o *option) *Ctx {
	p := &Ctx{
		o: o,

		createdAt: time.Now().Truncate(time.Second),
		context:   context,
		cache:     NewCache(),
	}

	var b bytes.Buffer
	b.WriteString("method:")
	b.WriteString(string(p.Method()))
	b.WriteString(" path:")
	b.WriteString(p.Path())

	b.WriteString(" ")
	b.WriteString(p.RealIp())
	b.WriteString(" ")

	b.WriteString(" header:")

	for k, v := range p.GetReqHeaderAll() {
		b.WriteString(k)
		b.WriteString(":")
		b.WriteString(v)
		b.WriteString(" ")
	}

	b.WriteString(" req:")
	b.Write(p.Body())

	log.Info("request ", b.String())

	return p
}

func (p *Ctx) Context() *fasthttp.RequestCtx {
	return p.context
}

func (p *Ctx) Method() Method {
	return Method(p.context.Request.Header.Method())
}

func (p *Ctx) Path() string {
	return string(p.context.Path())
}

func (p *Ctx) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Ctx) Close() {
	p.cache.Close()

	var b bytes.Buffer
	b.WriteString("method:")
	b.WriteString(string(p.Method()))
	b.WriteString(" path:")
	b.WriteString(p.Path())

	b.WriteString(" ")
	b.WriteString(p.RealIp())
	b.WriteString(" ")

	b.WriteString(" req:")
	b.Write(p.Body())

	b.WriteString(" rsp:")
	if p.respBody != nil {
		b.WriteString(*p.respBody)
		p.Context().Response.AppendBodyString(*p.respBody)
	}

	b.WriteString(" used:")
	b.WriteString(time.Since(p.CreatedAt()).String())

	log.Info("response ", b.String())
}

func (p *Ctx) RealIp() string {
	val := p.GetReqHeader("Cf-Connecting-Ip")
	if val != "" {
		if !netip.MustParseAddr(val).IsPrivate() {
			return val
		}
	}

	val = p.GetReqHeader("Fastly-Client-Ip")
	if val != "" {
		if !netip.MustParseAddr(val).IsPrivate() {
			return val
		}
	}

	val = p.GetReqHeader("True-Client-Ip")
	if val != "" {
		if !netip.MustParseAddr(val).IsPrivate() {
			return val
		}
	}

	val = p.GetReqHeader("X-Real-IP")
	if val != "" {
		if !netip.MustParseAddr(val).IsPrivate() {
			return val
		}
	}

	val = p.GetReqHeader("X-Client-IP")
	if val != "" {
		if !netip.MustParseAddr(val).IsPrivate() {
			return val
		}
	}

	val = p.GetReqHeader("X-Original-Forwarded-For")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if !netip.MustParseAddr(val).IsPrivate() {
				return v
			}
		}
		if !netip.MustParseAddr(val).IsPrivate() {
			return val
		}
	}

	val = p.GetReqHeader("X-Forwarded-For")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	val = p.GetReqHeader("X-Forwarded")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	val = p.GetReqHeader("Forwarded-For")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	val = p.GetReqHeader("Forwarded")
	if val != "" {
		for _, v := range strings.Split(val, ",") {
			if net.ParseIP(v) != nil {
				return v
			}
		}
		if net.ParseIP(val) != nil {
			return val
		}
	}

	return p.Context().RemoteIP().String()
}
