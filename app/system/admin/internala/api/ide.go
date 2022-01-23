package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internala/service"
	"scnu-coding/library/response"
)

var IDE = iDEApi{}

type iDEApi struct{}

func (a *iDEApi) ListAllIDE(r *ghttp.Request) {
	resp, err := service.IDE.ListContainer(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (a *iDEApi) RemoveIDE(r *ghttp.Request) {
	containerId := r.GetString("containerId")
	if err := service.IDE.RemoveContainer(r.Context(), containerId); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (a *iDEApi) RestartIDE(r *ghttp.Request) {
	containerId := r.GetString("containerId")
	if err := service.IDE.RestartContainer(r.Context(), containerId); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (a *iDEApi) GetServerInfo(r *ghttp.Request) {
	info, err := service.IDE.GetServerInfo(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, info)
}
