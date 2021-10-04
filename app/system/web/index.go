package web

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internal/api"
)

// @Author: 陈健航
// @Date: 2021/5/25 10:10
// @Description:

// Init 初始化路由
func Init() {
	s := g.Server()
	s.Group("/web", func(group *ghttp.RouterGroup) {
		group.Group("/user", func(group *ghttp.RouterGroup) {
			// 获取自己的用户信息
			group.GET("/myself", api.User.GetUserInfo)
			group.Group("/is-used", func(group *ghttp.RouterGroup) {
				// 邮箱是否已经被使用
				group.GET("/email", api.User.IsEmailUsed)
				// 学号是否已经被使用
				group.GET("/userNum", api.User.IsUserNumUsed)
			})
			// 更新个人信息
			group.POST("/", api.User.UpdateUserInfo)
		})
		group.Group("/ide", func(group *ghttp.RouterGroup) {
			group.POST("/", api.Ide.OpenIDE)
			group.Group("/alive", func(group *ghttp.RouterGroup) {
				group.POST("/", api.Ide.OpenFront)
				group.DELETE("/", api.Ide.CloseFront)
			})
		})
		group.Group("/checkin", func(group *ghttp.RouterGroup) {
			group.GET("/student", api.Checkin.StuListCheckinRecords)
			group.GET("/status", api.Checkin.GetCheckinStatus)
			group.Group("/record", func(group *ghttp.RouterGroup) {
				group.GET("/", api.Checkin.ListCheckinRecordByCourseId)
				group.POST("/", api.Checkin.StartCheckIn)
				group.DELETE("/", api.Checkin.DeleteCheckinRecord)
			})
			group.Group("/detail", func(group *ghttp.RouterGroup) {
				group.PUT("/", api.Checkin.CheckIn)
				group.GET("/", api.Checkin.ListCheckinDetailByCheckInRecordId)
			})
			group.GET("/export", api.Checkin.ExportCheckinRecord)
		})
		group.Group("/course", func(group *ghttp.RouterGroup) {
			group.GET("/is-enroll", api.Course.IsCourseEnroll)
			group.GET("/enroll", api.Course.ListCourseEnroll)
			group.GET("/teacher", api.Course.ListCourseByTeacherId)
			group.GET("/:courseId", api.Course.GetCourseDetail)
			group.GET("/overview", api.Course.ListCourseStudentOverview)
			group.GET("/search", api.Course.SearchCourseByCourseNameOrTeacherName)
			group.POST("/", api.Course.CreateCourse)
			group.Group("/announcement", func(group *ghttp.RouterGroup) {
				group.POST("/", api.CourseAnnouncement.InsertCourseAnnouncement)
				group.PUT("/", api.CourseAnnouncement.UpdateCourseAnnouncement)
				group.GET("/", api.CourseAnnouncement.ListCourseAnnouncement)
				group.DELETE("/", api.CourseAnnouncement.DeleteCourseAnnouncement)
			})
			group.GET("/open", api.Course.IsOpenByTeacherId)
			group.Group("/student", func(group *ghttp.RouterGroup) {
				group.POST("/", api.Course.ImportStudent2Class)
			})
			group.PUT("/", api.Course.UpdateCourse)
		})
		group.Group("/lab", func(group *ghttp.RouterGroup) {
			group.GET("/", api.Lab.ListLabByCourseId)
			group.DELETE("/", api.Lab.DeleteLab)
			group.PUT("/", api.Lab.UpdateLab)
			group.POST("/", api.Lab.InsertLab)
		})
		group.Group("/comment", func(group *ghttp.RouterGroup) {
			group.Group("/course", func(group *ghttp.RouterGroup) {
				group.GET("/", api.Comment.ListCourseComment)
				group.POST("/", api.Comment.InsertCourseComment)
				group.DELETE("/", api.Comment.DeleteCourseComment)
			})
			group.Group("/lab", func(group *ghttp.RouterGroup) {
				group.GET("/", api.Comment.ListLabComment)
				group.POST("/", api.Comment.InsertLabComment)
				group.DELETE("/", api.Comment.DeleteLabComment)
			})
		})
		group.Group("/submit", func(group *ghttp.RouterGroup) {
			group.Group("/", func(group *ghttp.RouterGroup) {
				group.GET("/", api.LabSummit.ListLabSubmit)
				group.POST("/", api.LabSummit.UpdateScoreAndComment)
			})
			group.PUT("/finish", api.LabSummit.UpdateFinishStat)
			group.Group("/report", func(group *ghttp.RouterGroup) {
				group.GET("/", api.LabSummit.GetReportContent)
				group.PUT("/", api.LabSummit.UpdateReportContent)
			})
			group.GET("/code", api.LabSummit.GetCode)
			group.GET("/id", api.LabSummit.ListLabSubmitId)
			group.POST("/correct", api.LabSummit.UpdateScoreAndComment)
			group.GET("/export", api.LabSummit.ExportScore)
		})
	})
}
