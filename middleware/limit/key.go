package limit

import (
	"bytes"

	"github.com/xihui-forever/goon"
)

func Ip(ctx *goon.Ctx) string {
	return ctx.RealIp()
}

func Path(ctx *goon.Ctx) string {
	return ctx.Path()
}

func MethodPath(ctx *goon.Ctx) string {
	var b bytes.Buffer
	b.WriteString(ctx.Method().String())
	b.WriteString("_")
	b.WriteString(ctx.Path())

	return b.String()
}

func MethodPathIp(ctx *goon.Ctx) string {
	var b bytes.Buffer
	b.WriteString(ctx.Method().String())
	b.WriteString("_")
	b.WriteString(ctx.Path())
	b.WriteString("_")
	b.WriteString(ctx.RealIp())

	return b.String()
}

func Header(key string) func(ctx *goon.Ctx) string {
	return func(ctx *goon.Ctx) string {
		return ctx.GetReqHeader(key)
	}
}
