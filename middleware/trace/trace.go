package trace

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
)

type Option struct {
	// 读trace id的header，如果不存在，则使用 `X-Goon-Trace` 作为header
	Header string `json:"header,omitempty"`

	// 生成trace id的函数
	GenTraceId func() string `json:"trace_id_func,omitempty"`
}

func Handler(opt Option) func(ctx *goon.Ctx) error {
	if opt.Header == "" {
		opt.Header = "X-Goon-Trace"
	}

	if opt.GenTraceId == nil {
		opt.GenTraceId = log.GenTraceId
	}

	return func(ctx *goon.Ctx) error {
		traceId := ctx.GetReqHeader(opt.Header)
		if traceId == "" {
			traceId = opt.GenTraceId()
		}
		log.SetTrace(traceId)
		defer log.DelTrace()

		return ctx.Next()
	}
}
