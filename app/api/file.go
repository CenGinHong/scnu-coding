package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/service"
	"scnu-coding/library/response"
	"time"
)

// @Author: 陈健航
// @Date: 2021/2/26 19:42
// @Description:

var File = fileApi{}

type fileApi struct{}

func (f *fileApi) UploadFile(r *ghttp.Request) {
	uploadFile := r.GetUploadFile("file")
	url, err := service.File.UploadFileAndGetUrl(r.Context(), uploadFile)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{"url": url})
}

func (f *fileApi) RemoveFile(r *ghttp.Request) {
	removeFile := r.GetString("file")
	if err := service.File.RemoveObject(r.Context(), removeFile); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (f *fileApi) GetObjectUrl(r *ghttp.Request) {
	filename := r.GetString("filename")
	url := service.File.GetObjectUrl(r.Context(), filename)
	response.Succ(r, g.Map{
		"url": url,
	})
}

func (f *fileApi) GetObjectPresignedUrl(r *ghttp.Request) {
	filename := r.GetString("filename")
	expire := r.GetInt("expire")
	url, err := service.File.GetObjectPresignedUrl(r.Context(), filename, time.Duration(expire)*time.Second)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{
		"url": url,
	})
}
