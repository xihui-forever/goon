package session

import (
	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/cache"
	"github.com/darabuchi/utils/xtime"
	"github.com/garyburd/redigo/redis"
	"github.com/nats-io/nuid"
)

var (
	nid = nuid.New()
)

func GenSession(data any) (string, error) {
	var sid string
	for i := 0; i < 3; i++ {
		sid = "session_" + nid.Next()
		ok, err := cache.SetNxWithTimeout(sid, data, xtime.Day)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}

		if ok {
			return sid, nil
		}
	}

	return "", ErrSessionGenerateFail
}

func GetSession(sid string) (string, error) {
	data, err := cache.Get(sid)
	if err != nil {
		if err == redis.ErrNil {
			return "", ErrSessionNotExist
		}
		log.Errorf("err:%v", err)
		return "", err
	}

	return data, nil
}

func GetSessionJson(sid string, obj any) error {
	data, err := GetSession(sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// sonic -> json
	err = sonic.UnmarshalString(data, obj)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}
