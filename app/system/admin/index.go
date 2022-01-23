package admin

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internala/api"
)

func Init() {
	s := g.Server()
	s.BindMiddlewareDefault()
	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.GET("/student", api.User.GetAllStudent)
			group.POST("/password", api.User.ResetPassword)
			group.PATCH("/", api.User.UpdateUser)
		})
		group.Group("/course", func(group *ghttp.RouterGroup) {
			group.GET("/", api.Course.ListAllCourse)
			group.Group("/enroll", func(group *ghttp.RouterGroup) {
				group.GET("/", api.Course.ListCourseEnroll)
				group.DELETE("/", api.Course.RemoveCourseEnroll)
			})
		})
		group.Group("/lab", func(group *ghttp.RouterGroup) {
			group.GET("/", api.Lab.ListAllLab)
		})
		group.Group("/submit", func(group *ghttp.RouterGroup) {
			group.GET("/", api.LabSubmit.ListAllSubmitByLabId)
		})
		group.Group("/ide", func(group *ghttp.RouterGroup) {
			group.GET("/", api.IDE.ListAllIDE)
			group.PUT("/", api.IDE.RestartIDE)
			group.DELETE("/", api.IDE.RemoveIDE)
			group.GET("/server", api.IDE.GetServerInfo)
		})
	})
}
