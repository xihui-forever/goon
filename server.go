package main

import (
	"errors"
	"github.com/darabuchi/log"
	"github.com/valyala/fasthttp"
	"github.com/xihui-forever/goon/handler"
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

	s1 := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			err := mux.CallOneOff(&ctx.Response, &ctx.Request)
			if err != nil {
				log.Errorf("err:%v", err)
				ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
				_, e := ctx.Write([]byte(err.Error()))
				if e != nil {
					log.Errorf("err:%v", e)
				}
			}
		},
	}

	s2 := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			err := mux.CallChrunked(&ctx.Response, &ctx.Request)
			if err != nil {
				log.Errorf("err:%v", err)
				ctx.Response.Header.SetStatusCode(fasthttp.StatusInternalServerError)
				_, e := ctx.Write([]byte(err.Error()))
				if e != nil {
					log.Errorf("err:%v", e)
				}
			}
		},
	}
	_ = s1.ListenAndServe(":8080")
	_ = s2.ListenAndServe(":8080")
}
