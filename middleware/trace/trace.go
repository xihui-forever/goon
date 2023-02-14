package trace

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
)

type Option struct {
	// 读trace id的header，如果不存在，则使用 `X-Goon-Trace` 作为header
	Header string `json:"header,omitempty"`
}

func Handler(opt Option) func(ctx *goon.Ctx) error {
	if opt.Header == "" {
		opt.Header = "X-Goon-Trace"
	}

	return func(ctx *goon.Ctx) error {
		traceId := ctx.GetReqHeader(opt.Header)
		if traceId == "" {
			traceId = log.GenTraceId()
		}
		log.SetTrace(traceId)
		defer log.DelTrace()

		return ctx.Next()
	}
}
