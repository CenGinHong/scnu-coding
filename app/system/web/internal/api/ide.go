package api

// @Author: 陈健航
// @Date: 2021/3/5 20:07
// @Description:

import (
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/system/web/internal/service"
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
	}
	url, err := service.Ide.OpenIDE(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, g.Map{"url": url})
}

// OpenFront
// @receiver i *ideAPI
// @param r *ghttp.Request
// @date 2021-07-17 22:28:18
func (i *ideAPI) OpenFront(r *ghttp.Request) {
	Id := r.GetInt("userId")
	labId := r.GetInt("labId")
	if err := service.Ide.OpenFront(r.Context(), 0, Id, labId); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, true)
}
