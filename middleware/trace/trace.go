package trace

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
)

const DefaultHeader = "X-Goon-Trace"

type Option struct {
	// 读trace id的header，如果不存在，则使用 `X-Goon-Trace` 作为header
	Header string `json:"header,omitempty"`

	// 生成trace id的函数
	GenTraceId func() string `json:"trace_id_func,omitempty"`
}

func Handler(opt Option) goon.Handler {
	if opt.Header == "" {
		opt.Header = DefaultHeader
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
