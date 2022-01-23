package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internala/service"
	"scnu-coding/library/response"
)

var LabSubmit = labSubmitApi{}

type labSubmitApi struct{}

func (a *labSubmitApi) ListAllSubmitByLabId(r *ghttp.Request) {
	labId := r.GetInt("labId")
	resp, err := service.LabSubmit.ListLabSubmitByLabId(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}
