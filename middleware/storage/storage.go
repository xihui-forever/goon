package storage

import (
	"time"
)

type Storage interface {
	Set(key string, value interface{}) error
	SetNx(key string, value interface{}) (bool, error)
	Get(key string) (string, error)
	Expire(key string, timeout time.Duration) error

	Inc(key string) (int64, error)
	IncBy(key string, value int64) (int64, error)
	Dec(key string) (int64, error)
	DecBy(key string, value int64) (int64, error)

	Close() error
}
