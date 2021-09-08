package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internal/service"
	"scnu-coding/library/response"
)

var Course = courseApi{}

type courseApi struct{}

func (a *courseApi) GetAllCourse(r *ghttp.Request) {
	resp, err := service.Course.GetAllCourse(r.Context())
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}
