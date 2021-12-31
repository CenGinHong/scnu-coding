package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internala/service"
	"scnu-coding/library/response"
)

var Course = courseApi{}

type courseApi struct{}

func (a *courseApi) ListAllCourse(r *ghttp.Request) {
	resp, err := service.Course.ListAllCourse(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (a courseApi) ListCourseEnroll(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Course.ListEnroll(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}
