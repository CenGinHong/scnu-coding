package admin

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/admin/internal/api"
)

func Init() {
	s := g.Server()
	s.BindMiddlewareDefault()
	s.Group("/admin", func(group *ghttp.RouterGroup) {
		group.Group("/user", func(group *ghttp.RouterGroup) {
			group.GET("/", api.User.GetAllUser)
			group.POST("/password", api.User.ResetPassword)
			group.PATCH("/", api.User.UpdateUser)
		})
		group.Group("/course", func(group *ghttp.RouterGroup) {
			group.GET("/", api.Course.ListAllCourse)
			group.Group("/enroll", func(group *ghttp.RouterGroup) {
				group.GET("/", api.Course.ListCourseEnroll)
			})
		})
	})
}
