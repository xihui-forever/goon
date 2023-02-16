package limit

import (
	"time"
)

type Storage interface {
	Inc(key string) (int64, error)
	Dec(key string) (int64, error)
	DecBy(key string, value int64) (int64, error)
	Expire(key string, timeout time.Duration) error
}

type Config struct {
	Key        string
	Expiration time.Duration
	Storage    Storage
	Max        int64
}
