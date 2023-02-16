package limit

import (
	"time"
)

type Storage interface {
	Inc(key string) (int64, error)
	Dec(key string) (int64, error)
	DecBy(key string, value int64) (int64, error)
	Expire(key string, timeout time.Duration) error
	AddItem(key string, value int64) error
	GetNotValid(key string, timeout time.Duration) (int64, error)
	DeleteItem(key string, value int64) error
	LenItemList(key string) (int64, error)
}

type Config struct {
	Key        string
	Expiration time.Duration
	Storage    Storage
	Max        int64
}
