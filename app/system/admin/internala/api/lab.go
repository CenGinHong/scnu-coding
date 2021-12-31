package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internala/service"
	"scnu-coding/library/response"
)

var Lab = labApi{}

type labApi struct{}

func (a *labApi) ListAllLab(r *ghttp.Request) {
	resp, err := service.Lab.ListAllLab(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}
