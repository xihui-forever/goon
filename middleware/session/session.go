package session

import (
	"time"

	"github.com/bytedance/sonic"
	"github.com/darabuchi/log"
	"github.com/nats-io/nuid"
	"github.com/xihui-forever/goon/ctx"
	"github.com/xihui-forever/goon/middleware/storage"
	"github.com/xihui-forever/goon/middleware/storage/memory"
)

var (
	nid = nuid.New()

	def = New(memory.New(memory.Option{
		CleanType: memory.CleanFasterObject,
		MaxSize:   1024,
		DataType:  memory.Lru,
	}))
)

type Session struct {
	storage Storage
}

func New(storage Storage) *Session {
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

type Option[M any] struct {
	// 如果不存在，这使用默认的session
	Session *Session `json:"session,omitempty"`

	// 读写Session的header，如果不存在，则使用 `X-Goon-Session` 作为header
	Header string `json:"header,omitempty"`

	// 过期时间，如果不存在，则使用默认的过期时间，time.Hour
	Expiration time.Duration `json:"expiration,omitempty"`

	OnError   func(ctx *ctx.Ctx, err error) error `json:"on_error,omitempty"`
	OnSuccess func(ctx *ctx.Ctx, obj M) error     `json:"on_success,omitempty"`
}

func Handler[M any](opt Option[M]) func(ctx *ctx.Ctx) error {
	if opt.Session == nil {
		opt.Session = def
	}

	if opt.Header == "" {
		opt.Header = "X-Goon-Session"
	}

	return func(ctx *ctx.Ctx) error {
		sid := ctx.GetReqHeader(opt.Header)

		data, err := opt.Session.GetSession(sid)
		if err != nil {
			log.Errorf("err:%v", err)
			if opt.OnError != nil {
				return opt.OnError(ctx, err)
			}
			return err
		}

		var obj M
		err = sonic.UnmarshalString(data, &obj)
		if err != nil {
			log.Errorf("err:%v", err)
			if opt.OnError != nil {
				return opt.OnError(ctx, err)
			}
			return err
		}

		if opt.OnSuccess != nil {
			return opt.OnSuccess(ctx, obj)
		}

		return nil
	}
}
