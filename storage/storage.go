package storage

import (
	"time"
)

type Storage interface {
	Connect(addr string, db int, password string) error
	SetNx(key string, value interface{}) (bool, error)
	Get(key string) (string, error)
	Expire(key string, timeout time.Duration) error
}
