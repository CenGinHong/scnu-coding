package api

import (
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/system/web/internala/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/3/1 0:32
// @Description:

var CourseAnnouncement = courseAnnouncementAPI{}

type courseAnnouncementAPI struct{}

// Insert 插入课程资源
// @receiver receiver
// @params r
// @date 2021-05-25 00:12:51
func (c *courseAnnouncementAPI) Insert(r *ghttp.Request) {
	var req *define.InsertCourseAnnouncementReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.CourseAnnouncement.InsertCourseAnnouncement(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, true)
}

func (c *courseAnnouncementAPI) Update(r *ghttp.Request) {
	var req *define.UpdateCourseAnnouncementReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.CourseAnnouncement.UpdateCourseAnnouncement(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (c *courseAnnouncementAPI) ListByCourseId(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.CourseAnnouncement.ListCourseAnnouncement(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (c *courseAnnouncementAPI) Delete(r *ghttp.Request) {
	courseAnnouncementId := r.GetInt("id")
	if err := service.CourseAnnouncement.DeleteCourseResource(r.Context(), courseAnnouncementId); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, true)
}

func (c *courseAnnouncementAPI) GetOne(r *ghttp.Request) {
	courseAnnouncementId := r.GetInt("id")
	resp, err := service.CourseAnnouncement.GetOne(r.Context(), courseAnnouncementId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}
