package session

import (
	"time"
)

type Storage interface {
	SetNx(key string, value interface{}) (bool, error)
	Get(key string) (string, error)
	Expire(key string, timeout time.Duration) error
	Close() error
}
