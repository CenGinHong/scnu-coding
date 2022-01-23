package api

// @Author: 陈健航
// @Date: 2021/3/5 20:07
// @Description:

import (
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/system/web/internala/service"
	"scnu-coding/library/response"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

// Ide Ide对外API
var Ide = ideAPI{}

type ideAPI struct{}

// OpenIDE 打开容器
// @receiver receiver
// @params r
// @date 2021-05-25 10:28:44
func (i *ideAPI) OpenIDE(r *ghttp.Request) {
	var req *define.OpenIDEReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	url, err := service.IDE.OpenIDE(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{"url": url})
}

// FrontAlive ide前端插件发来的存活信息
// @receiver i *ideAPI
// @param r *ghttp.Request
// @date 2021-07-17 22:28:18
func (i *ideAPI) FrontAlive(r *ghttp.Request) {
	var req *define.FrontAliveReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	if err := service.IDE.FrontAlive(r.Context(), req); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r)
}
