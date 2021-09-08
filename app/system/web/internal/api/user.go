package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/system/web/internal/service"
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
	}
	response.Succ(r, resp)
}

func (u *userApi) IsEmailUsed(r *ghttp.Request) {
	email := r.GetString("email")
	used, err := service.SysUser.IsEmailUsed(r.Context(), email)
	if err != nil {
		response.Exit(r, err)
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
	}
	if err := service.SysUser.Update(r.Context(), req); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r)
}

func (u *userApi) ListCodingTimeByUserId(r *ghttp.Request) {
	var req *define.ListCodingTimeByUserIdReq
	service.SysUser.ListCodingTimeByUserId(r.Context(), req)

}
