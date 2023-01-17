package main

import (
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils/cache"
	"github.com/darabuchi/utils/db"
	"github.com/darabuchi/utils/xtime"
	"github.com/spf13/viper"
	"github.com/xihui-forever/goon/admin"
	"github.com/xihui-forever/goon/config"
	"github.com/xihui-forever/goon/handler"
	"github.com/xihui-forever/goon/session"
	"github.com/xihui-forever/goon/types"
	"net/http"
	"time"
)

func main() {
	config.Load()

	err := cache.Connect("127.0.0.1:6379", 1, "")
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	err = db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelAdmin{},
	)
	if err != nil {
		log.Errorf("err:%v", err)
		return
	}

	mux := handler.NewHandler()

	mux.Post("/admin/login", func(ctx *handler.Ctx, req *types.AdminLoginReq) (*types.AdminLoginRsp, error) {
		var rsp types.AdminLoginRsp

		a, err := admin.GetAdmin(req.Username)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		if a.Password != admin.Encrypt(req.Password) {
			return nil, admin.ErrAdminExist
		}

		rsp.Token, err = session.GenSession(&types.AdminSession{
			Id:       a.Id,
			Username: a.Username,
		})
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		rsp.Expire = uint32(time.Now().Add(xtime.Day).Unix())

		return &rsp, nil
	})

	mux.PreUse("/user", func(ctx *handler.Ctx) error {
		var sess types.AdminSession
		err := session.GetSessionJson(ctx.GetSid(), &sess)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		log.Infof("sess:%v", sess)

		return nil
	})

	mux.Post("/user/get", func(ctx *handler.Ctx) error {
		return nil
	})

	_ = http.ListenAndServe(":8080", mux)
}
