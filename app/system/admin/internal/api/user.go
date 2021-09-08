package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internal/define"
	"scnu-coding/app/system/admin/internal/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/4/29 14:32
// @Description:

var User = userApi{}

type userApi struct{}

func (a *userApi) GetAllUser(r *ghttp.Request) {
	resp, err := service.SysUser.GetAllUser(r.Context())
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

func (a *userApi) ResetPassword(r *ghttp.Request) {
	userId := r.GetInt("userId")
	if err := service.SysUser.ResetPassword(r.Context(), userId); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r)
}

func (a *userApi) UpdateUser(r *ghttp.Request) {
	var req define.UpdateSysUserReq
	err := service.SysUser.UpdateUser(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r)
}
