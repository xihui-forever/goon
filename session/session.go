package session

import (
	"time"

	"github.com/darabuchi/log"
	"github.com/nats-io/nuid"
	"github.com/xihui-forever/goon/storage"
	"github.com/xihui-forever/goon/storage/memory/hashmap"
)

var (
	nid = nuid.New()

	def = NewSession(hashmap.New())
)

type Session struct {
	storage storage.Storage
}

func NewSession(storage storage.Storage) *Session {
	return &Session{
		storage: storage,
	}
}

func (s *Session) GenSession(data any, timeout time.Duration) (string, error) {
	var sid string
	for i := 0; i < 3; i++ {
		sid = "session_" + nid.Next()
		ok, err := s.storage.SetNx(sid, data)
		if err != nil {
			log.Errorf("err:%v", err)
			return "", err
		}
		if ok {
			err = s.storage.Expire(sid, timeout)
			if err != nil {
				log.Errorf("err:%v", err)
			}
			return sid, nil
		}
	}

	return "", ErrSessionGenerateFail
}

func (s *Session) GetSession(sid string) (string, error) {
	data, err := s.storage.Get(sid)
	if err != nil {
		if err == storage.ErrKeyNotExist {
			return "", ErrSessionNotExist
		}

		log.Errorf("err:%v", err)
		return "", err
	}

	return data, nil
}

func (s *Session) Expire(sid string, t time.Duration) error {
	err := s.storage.Expire(sid, t)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func (s *Session) SetStorage(storage storage.Storage) {
	s.storage = storage
}

func GenSession(data any, timeout time.Duration) (string, error) {
	return def.GenSession(data, timeout)
}

func GetSession(sid string) (string, error) {
	return def.GetSession(sid)
}

func Expire(sid string, t time.Duration, storage storage.Storage) error {
	return def.Expire(sid, t)
}

func SetStorage(s storage.Storage) {
	def.SetStorage(s)
}
