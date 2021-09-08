package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/2/26 19:42
// @Description:

var File = fileApi{}

type fileApi struct{}

func (f *fileApi) UploadFile(r *ghttp.Request) {
	uploadFile := r.GetUploadFile("file")
	url, err := service.File.UploadFile(r.Context(), uploadFile)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, g.Map{"url": url})
}
