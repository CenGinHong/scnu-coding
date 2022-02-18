package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internala/define"
	"scnu-coding/app/system/admin/internala/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/4/29 14:32
// @Description:

var User = userApi{}

type userApi struct{}

func (a *userApi) ListUser(r *ghttp.Request) {
	resp, err := service.SysUser.ListUser(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (a *userApi) GetUser(r *ghttp.Request) {
	userId := r.GetInt("id")
	resp, err := service.SysUser.GetUser(r.Context(), userId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (a *userApi) ResetPassword(r *ghttp.Request) {
	userId := r.GetInt("id")
	if err := service.SysUser.ResetPassword(r.Context(), userId); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (a *userApi) UpdateUser(r *ghttp.Request) {
	var req define.UpdateSysUserReq
	err := service.SysUser.UpdateUser(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (a *userApi) GetImportDemoCsv(r *ghttp.Request) {
	csv, err := service.SysUser.GetImportDemoCsv(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	r.Response.Header().Set("Content-Disposition", "attachment;filename=demo.xlsx")
	// 响应类型,编码
	r.Response.Header().Set("Content-Type", "application/vnd.ms-excel;charset=utf8")
	r.Response.WriteExit(csv)
}

func (a *userApi) ImportUserIdByCsv(r *ghttp.Request) {
	file := r.GetUploadFile("file")
	roleId := r.GetInt("roleId")
	errMsg, err := service.SysUser.ImportStudent(r.Context(), file, roleId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, errMsg)
}
func (a *userApi) DeleteUser(r *ghttp.Request) {
	userId := r.GetInt("id")
	if err := service.SysUser.DeleteUser(r.Context(), userId); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}
