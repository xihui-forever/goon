package session

import (
	"github.com/darabuchi/log"
	"github.com/nats-io/nuid"
	"github.com/xihui-forever/goon/handler"
	"github.com/xihui-forever/goon/storage"
	"time"
)

var (
	nid = nuid.New()
)

func HasSession(ctx *handler.Ctx, storage storage.Storage) (bool, error) {
	sid := ctx.GetSid()
	_, err := storage.Get(sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return false, err
	}
	return true, nil
}

func GenSession(data any, storage storage.Storage) (err error) {
	var sid string
	for i := 0; i < 3; i++ {
		sid = "session_" + nid.Next()
		ok, err := storage.SetNx(sid, data)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}
		if ok {
			return nil
		}
	}
	return ErrSessionGenerateFail
}

func GetSession(sid string, storage storage.Storage) (string, error) {
	data, err := storage.Get(sid)
	if err != nil {
		log.Errorf("err:%v", err)
		return "", err
	}
	return data, nil
}

func SetSessionExpire(sid string, t time.Duration, storage storage.Storage) error {
	return storage.Expire(sid, t)
}
