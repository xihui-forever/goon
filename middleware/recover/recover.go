package recover

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
)

type Option struct {
	Next func(ctx *goon.Ctx) error
}

func Handler(opt Option) goon.Handler {
	return func(ctx *goon.Ctx) (err error) {
		defer func() {
			// 接收panic的信息，防止某一个请求导致程序崩溃
			if e := recover(); e != nil {
				log.Errorf("PANIC err:%v", e)
				if opt.Next != nil {
					err = opt.Next(ctx)
				}
			}
		}()

		return ctx.Next()
	}
}
