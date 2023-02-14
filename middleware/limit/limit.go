package limit

import (
	"errors"
	"github.com/xihui-forever/goon/middleware/storage/memory/mem"
	"sync"
	"time"

	"github.com/xihui-forever/goon"
)

type Option struct {
	// 是否要跳过限频
	NeedSkip func(ctx *goon.Ctx) bool

	// 限频的唯一标识，如果为空，则使用 Ip 作为唯一标识
	KeyGenerator func(ctx *goon.Ctx) string

	// 如果为空，则使用内存存储
	Storage Storage

	// 限频的过期时间，针对不同的 LimiterMiddleware 含义不同
	//
	// Default: 1 * time.Minute
	Expiration time.Duration

	// 指定时间范围内的最大请求数，如果不存在，则为 10
	Max int64

	// 限频的逻辑，如果为空，则使用 FixedWindow
	LimiterMiddleware func(cfg Config) bool
}

func Handler(opt Option) goon.Handler {
	if opt.KeyGenerator == nil {
		opt.KeyGenerator = Ip
	}

	if opt.Expiration == 0 {
		opt.Expiration = 1 * time.Minute
	}

	if opt.Storage == nil {
		opt.Storage = mem.New(true)
	}

	if opt.LimiterMiddleware == nil {
		opt.LimiterMiddleware = FixedWindow
	}

	pool := sync.Pool{
		New: func() interface{} {
			return &Config{
				Expiration: opt.Expiration,
				Storage:    opt.Storage,
				Max:        opt.Max,
			}
		},
	}

	return func(ctx *goon.Ctx) error {
		if opt.NeedSkip != nil && opt.NeedSkip(ctx) {
			return ctx.Next()
		}

		key := opt.KeyGenerator(ctx)

		cfg := pool.Get().(*Config)
		defer pool.Put(cfg)

		cfg.Key = key

		if opt.LimiterMiddleware(*cfg) {
			return errors.New("too many requests")
		}

		return ctx.Next()
	}
}
