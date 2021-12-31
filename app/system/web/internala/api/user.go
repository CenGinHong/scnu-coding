package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	service2 "scnu-coding/app/service"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/system/web/internala/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/4/29 14:32
// @Description:

var User = userApi{}

type userApi struct{}

func (u *userApi) GetUserInfo(r *ghttp.Request) {
	resp, err := service.SysUser.GetUserInfo(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (u *userApi) IsEmailUsed(r *ghttp.Request) {
	email := r.GetString("email")
	used, err := service.SysUser.IsEmailUsed(r.Context(), email)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{
		"isUsed": used,
	})
}

func (u *userApi) IsUserNumUsed(r *ghttp.Request) {
	userNum := r.GetString("userNum")
	used, err := service.SysUser.IsUserNumUsed(r.Context(), userNum)
	if err != nil {
		response.Exit(r, err)
		return
	}
	println(used)
	response.Succ(r, g.Map{
		"isUsed": used,
	})
}

func (u *userApi) UpdateUserInfo(r *ghttp.Request) {
	var req *define.UpdateUserInfoReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.SysUser.Update(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (u *userApi) ListCodingTimeByUserId(r *ghttp.Request) {
	var req *define.ListCodingTimeByUserIdReq
	err := r.Parse(&req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	// 如果没有指定userId则查自己的
	if req.UserId == 0 {
		req.UserId = service2.Context.Get(r.Context()).User.UserId
	}
	resp, err := service.SysUser.ListCodingTimeByUserId(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}
