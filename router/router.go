package router

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/api"
	"scnu-coding/app/service"
	"scnu-coding/app/system/admin"
	"scnu-coding/app/system/web"
	"scnu-coding/app/utils"
)

func init() {
	web.Init()
	admin.Init()
	s := g.Server()
	s.BindMiddlewareDefault(service.Middleware.Ctx)
	s.BindMiddlewareDefault(service.Middleware.CORS)
	s.Group("/", func(group *ghttp.RouterGroup) {
		utils.GfToken.Middleware(group)
		group.ALL("/hello/:id", api.Hello.Index)
	})
	s.Group("/test", func(group *ghttp.RouterGroup) {
		group.GET("/hello/:id", api.Hello.Index)
		group.GET("/", api.Hello.Index)
		group.POST("/upload", api.Hello.Index1)
		group.GET("/start", api.Hello.Index1)
	})
	s.Group("/file", func(group *ghttp.RouterGroup) {
		group.GET("/", api.File.GetObjectUrl)
		group.GET("/presigned", api.File.GetObjectPresignedUrl)
		group.POST("/", api.File.UploadFile)
		group.DELETE("/", api.File.RemoveFile)
	})
	s.Group("/common", func(group *ghttp.RouterGroup) {
		group.POST("/verCode", api.Common.SendVerCode)
	})
	s.BindHandler("/chat", api.Hello.Index1)
}
