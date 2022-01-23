package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internala/define"
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

func (a *courseApi) ListCourseEnroll(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Course.ListEnroll(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (a *courseApi) RemoveCourseEnroll(r *ghttp.Request) {
	var req *define.RemoveCourseEnrollReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.Course.RemoveCourseEnroll(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}
