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

func (a *courseApi) AddStudent2Class(r *ghttp.Request) {
	var req *define.AddStudent2ClassReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	errMsg, err := service.Course.AddStudent2Class(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, errMsg)
}
