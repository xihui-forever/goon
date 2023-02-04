package redis

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/cache"
	"github.com/garyburd/redigo/redis"
	"time"
)

type Redis struct{}

func (r Redis) Connect(addr string, db int, password string) error {
	err := cache.Connect(addr, db, password)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func (r Redis) SetNx(sid string, data interface{}) (bool, error) {
	ok, err := cache.SetNx(sid, data)
	if err != nil {
		log.Errorf("err:%v", err)
		return false, err
	}

	if ok {
		return true, nil
	}

	return ok, err
}

func (r Redis) Get(sid string) (string, error) {
	data, err := cache.Get(sid)
	if err != nil {
		if err == redis.ErrNil {
			return "", ErrSidNotExist
		}
		log.Errorf("err:%v", err)
		return "", err
	}

	return data, nil
}

func (r Redis) Expire(sid string, timeout time.Duration) error {
	ok, err := cache.Expire(sid, timeout)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	if !ok {
		return ErrRedisSetExpireFailure
	}
	return nil
}
