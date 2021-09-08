package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/system/web/internal/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/2/1 23:44
// @Description:

var Lab = labApi{}

type labApi struct{}

func (l *labApi) ListLabByCourseId(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Lab.ListLabByCourseId(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

//func (l *labApi) InsertLab(r *ghttp.Request) {
//	var req *define.InsertLabReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	if err := service.Lab.InsertLab(r.Context(), req); err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, true)
//}

func (l *labApi) UpdateLab(r *ghttp.Request) {
	var req *define.UpdateLabReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	if err := service.Lab.Update(r.Context(), req); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, true)
}

//func (receiver *labApi) GetOne(r *ghttp.Request) {
//	labId := r.GetInt("labId")
//	resp, err := service.Lab.GetOne(labId)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, resp)
//}

func (l *labApi) DeleteLab(r *ghttp.Request) {
	labId := r.GetInt("labId")
	// 查看开实验的人是不是用户
	if err := service.Lab.Delete(r.Context(), labId); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, true)
}

//
//func (receiver *labApi) ListByToken(r *ghttp.Request) {
//	var req *model.ListLabByTokenReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	req.StuID = r.GetCtxVar(dao.SysUser.Columns.StuID).Int()
//	resp, err := service.Lab.ListByToken(req)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, resp)
//}

//func (receiver *labApi) CheckCode(r *ghttp.Request) {
//	var req *model.CheckCodeReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	req.TeacherId = r.GetCtxVar(dao.SysUser.Columns.StuID).Int()
//	url, err := ide.Ide.CheckCode(req)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, g.Map{"url": url})
//}

//func (receiver labApi) ListLabScore(r *ghttp.Request) {
//	var req *model.ListLabScoreReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	if req.StuID == 0 {
//		req.StuID = r.GetCtxVar(dao.SysUser.Columns.StuID).Int()
//	}
//	resp, err := service.Lab.ListLabScore(req)
//	if err != nil {
//		return
//	}
//	response.Succ(r, resp)
//}
