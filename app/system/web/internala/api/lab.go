package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/system/web/internala/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/2/1 23:44
// @Description:

var Lab = labApi{}

type labApi struct{}

func (l *labApi) ListByCourseId(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Lab.ListLabByCourseId(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (l labApi) Insert(r *ghttp.Request) {
	var req *define.InsertLabReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	req.UploadFile = r.GetUploadFile("file")
	err := service.Lab.InsertLab(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (l *labApi) Update(r *ghttp.Request) {
	var req *define.UpdateLabReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	req.UploadFile = r.GetUploadFile("file")
	if err := service.Lab.Update(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, true)
}

func (l *labApi) Delete(r *ghttp.Request) {
	labId := r.GetInt("labId")
	// 查看开实验的人是不是用户
	if err := service.Lab.Delete(r.Context(), labId); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, true)
}

func (l *labApi) GetOne(r *ghttp.Request) {
	labId := r.GetInt("id")
	resp, err := service.Lab.GetOne(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}
