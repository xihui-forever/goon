package redis

import (
	"time"

	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/garyburd/redigo/redis"
	"github.com/shomali11/xredis"
	"github.com/xihui-forever/goon/middleware/storage"
)

type Redis struct {
	client *xredis.Client
}

func (r *Redis) Set(key string, value interface{}) error {
	panic("implement me")
}

func (r *Redis) Inc(key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) IncBy(key string, value int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) Dec(key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) DecBy(key string, value int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) ZAdd(key string, members ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) ZRange(key string, start, stop int64) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) ZRem(key string, members ...interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (r *Redis) ZLen(key string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func New(cfg *RedisConfig) *Redis {
	var opt []redis.DialOption
	if cfg.Password != "" {
		opt = append(opt, redis.DialPassword(cfg.Password))
	}
	if cfg.DB != 0 {
		opt = append(opt, redis.DialDatabase(cfg.DB))
	}

	opt = append(opt,
		redis.DialConnectTimeout(time.Second*3),
		redis.DialReadTimeout(time.Second*3),
		redis.DialWriteTimeout(time.Second*3),
		redis.DialKeepAlive(time.Minute),
	)

	return &Redis{
		client: xredis.NewClient(&redis.Pool{
			Dial: func() (redis.Conn, error) {
				return redis.Dial("tcp", cfg.Addr, opt...)
			},
			MaxIdle:     100,
			MaxActive:   100,
			IdleTimeout: time.Second * 5,
			Wait:        true,
		}),
	}
}

func (r *Redis) SetNx(key string, data interface{}) (bool, error) {
	ok, err := r.client.SetNx(key, utils.ToString(data))
	if err != nil {
		log.Errorf("err:%v", err)
		return false, err
	}

	if ok {
		return true, nil
	}

	return ok, err
}

func (r *Redis) Get(key string) (string, error) {
	val, ok, err := r.client.Get(key)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", storage.ErrKeyNotExist
	}

	return val, nil
}

func (r *Redis) Expire(key string, timeout time.Duration) error {
	_, err := r.client.Expire(key, int64(timeout.Seconds()))
	if err != nil {
		if err == redis.ErrNil {
			return storage.ErrKeyNotExist
		}

		return err
	}

	return nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}
