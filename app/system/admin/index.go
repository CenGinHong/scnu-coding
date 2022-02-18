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
			group.GET("/student", api.User.ListUser)
			group.PUT("/password/:id", api.User.ResetPassword)
			group.PUT("/", api.User.UpdateUser)
			group.DELETE("/:id", api.User.DeleteUser)
			group.Group("/import", func(group *ghttp.RouterGroup) {
				group.GET("/demo", api.User.GetImportDemoCsv)
				group.POST("/", api.User.ImportUserIdByCsv)
			})
			group.GET("/:id", api.User.GetUser)
		})
		group.Group("/course", func(group *ghttp.RouterGroup) {
			group.GET("/", api.Course.ListAllCourse)
			group.Group("/enroll", func(group *ghttp.RouterGroup) {
				group.GET("/", api.Course.ListCourseEnroll)
				group.PUT("/", api.Course.AddStudent2Class)
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
