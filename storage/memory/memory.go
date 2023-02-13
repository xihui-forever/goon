package memory

import (
	"time"

	"github.com/darabuchi/utils"
	"github.com/xihui-forever/goon/storage/memory/fifo"
	"github.com/xihui-forever/goon/storage/memory/hashmap"
	"github.com/xihui-forever/goon/storage/memory/lifo"
	"github.com/xihui-forever/goon/storage/memory/lru"
	"go.uber.org/atomic"
)

type MemoryInerface interface {
	SetNx(key string, value string) (bool, error)
	Get(key string) (string, error)
	Expire(key string, timeout time.Duration) error
	Clean(func(size int64) (stop bool))
}

type CleanType int

const (
	CleanFasterSize CleanType = iota
	CleanFasterObject
	CleanBetterSize
	CleanBetterObject
)

type DataType int

const (
	Hashmap DataType = iota
	Fifo
	Lru
	Lifo
)

type Option struct {
	CleanType CleanType

	MaxSize int64

	DataType DataType
}

type Memory struct {
	inter MemoryInerface
	opt   Option

	size *atomic.Int64
	c    chan struct{}
}

func (p *Memory) SetNx(key string, value interface{}) (ok bool, err error) {
	v := utils.ToString(value)
	ok, err = p.inter.SetNx(key, v)
	if err != nil {
		return
	}

	if ok {
		switch p.opt.CleanType {
		case CleanFasterSize:
			if p.size.Add(int64(len(v))) > p.opt.MaxSize {
				p.inter.Clean(func(size int64) (stop bool) {
					if p.size.Add(-size) <= p.opt.MaxSize {
						stop = true
					}
					return
				})
			}
		case CleanFasterObject:
			if p.size.Inc() > p.opt.MaxSize {
				p.inter.Clean(func(size int64) (stop bool) {
					if p.size.Dec() <= p.opt.MaxSize {
						stop = true
					}
					return
				})
			}
		case CleanBetterSize:
			p.size.Add(int64(len(v)))
		case CleanBetterObject:
			p.size.Inc()
		}
	}
	return
}

func (p *Memory) Get(key string) (string, error) {
	return p.inter.Get(key)
}

func (p *Memory) Expire(key string, timeout time.Duration) error {
	return p.inter.Expire(key, timeout)
}

func (p *Memory) Close() error {
	select {
	case <-p.c:
	}
	return nil
}

func NewMemory(opt Option) *Memory {
	p := &Memory{
		opt:  opt,
		size: atomic.NewInt64(0),

		c: make(chan struct{}, 1),
	}

	switch opt.DataType {
	case Hashmap:
		p.inter = hashmap.New()
	case Fifo:
		p.inter = fifo.New()
	case Lru:
		p.inter = lru.New()
	case Lifo:
		p.inter = lifo.New()
	default:
		panic("not support")
	}

	switch p.opt.CleanType {
	case CleanBetterSize, CleanBetterObject:
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					p.inter.Clean(func(size int64) (stop bool) {
						return false
					})
				case <-p.c:
					return
				}
			}
		}()
	}

	return p
}
