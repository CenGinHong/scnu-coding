package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/service"
	"scnu-coding/library/response"
)

var Common = commonApi{}

type commonApi struct{}

func (a *commonApi) SendVerCode(r *ghttp.Request) {
	email := r.GetString("email")
	if err := service.Common.SendVerificationCode(r.Context(), email); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}
