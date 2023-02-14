package recover

import (
	"github.com/darabuchi/log"
	"github.com/xihui-forever/goon"
)

type Option struct {
	OnPanic func(ctx *goon.Ctx, err interface{}) error
}

func Handler(opt Option) goon.Handler {
	return func(ctx *goon.Ctx) (err error) {
		defer func() {
			// 接收panic的信息，防止某一个请求导致程序崩溃
			if e := recover(); e != nil {
				log.Errorf("PANIC err:%v", e)
				if opt.OnPanic != nil {
					err = opt.OnPanic(ctx, e)
				}
			}
		}()

		return ctx.Next()
	}
}
