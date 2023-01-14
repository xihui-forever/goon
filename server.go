package main

import (
	"errors"
	"github.com/xihui-forever/goStudy/handler"
	"net/http"
)

type (
	GetUserReq struct {
		A string `json:"a"`
		B string `json:"b"`
	}

	GetUserRsp struct {
		A string `json:"a"`
		B string `json:"b"`
	}
)

type (
	SetUserReq struct {
		C string `json:"c"`
		D string `json:"d"`
	}

	SetUserRsp struct {
		C string `json:"c"`
		D string `json:"d"`
	}
)

func main() {
	mux := handler.NewHandler()

	mux.PreUse("/User", func(ctx *handler.Ctx) error {
		return nil
	})
	//group := group.use()

	mux.Head("/ ", func(ctx *handler.Ctx) error {
		return nil
	})

	mux.Get("/GetMe", func(ctx *handler.Ctx) (**GetUserRsp, error) {

		return nil, errors.New("not handle")
	})

	mux.Post("/SetMe", func(ctx *handler.Ctx, req *SetUserReq) error {

		return errors.New("not handle")
	})

	mux.Post("/User/SetUser", func(ctx *handler.Ctx, req *SetUserReq) error {

		return errors.New("not handle")
	})

	mux.PostUse("/User", func(ctx *handler.Ctx) error {

		return nil
	})

	_ = http.ListenAndServe(":8080", mux)
}
