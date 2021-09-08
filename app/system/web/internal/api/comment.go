package api

// @Author: 陈健航
// @Date: 2021/1/16 0:32
// @Description:

import (
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/system/web/internal/service"
	"scnu-coding/library/response"

	"github.com/gogf/gf/net/ghttp"
)

var Comment = commentAPI{}

type commentAPI struct{}

// InsertCourseComment 新增课程评论
// @receiver c
// @params r
// @date 2021-01-16 20:38:52
func (c *commentAPI) InsertCourseComment(r *ghttp.Request) {
	var req *define.InsertCourseCommentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	id, err := service.Comment.InsertCourseComment(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, g.Map{
		"id": id,
	})
}

// InsertLabComment 新增课程评论
// @receiver c
// @params r
// @date 2021-01-16 20:38:52
func (c *commentAPI) InsertLabComment(r *ghttp.Request) {
	var req *define.InsertLabCommentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	id, err := service.Comment.InsertLabComment(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, g.Map{
		"id": id,
	})
}

// ListCourseComment 分页查询课程评论
// @receiver c
// @params r
// @date 2021-01-30 21:43:14
func (c *commentAPI) ListCourseComment(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Comment.ListCourseCommentByCourseId(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

// ListLabComment 分页查询实验评论
// @receiver c
// @params r
// @date 2021-01-30 21:43:14
func (c *commentAPI) ListLabComment(r *ghttp.Request) {
	labId := r.GetInt("labId")
	resp, err := service.Comment.ListLabCommentByLabId(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

// DeleteCourseComment 删除课程评论,伪删除
// @receiver c
// @params r
// @date 2021-01-16 00:42:58
func (c *commentAPI) DeleteCourseComment(r *ghttp.Request) {
	courseCommentId := r.GetInt("courseCommentId")
	if err := service.Comment.DeleteCourseComment(r.Context(), courseCommentId); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, true)
}

// DeleteLabComment 删除实验评论,伪删除
// @receiver c
// @params r
// @date 2021-01-16 00:42:58
func (c *commentAPI) DeleteLabComment(r *ghttp.Request) {
	labCommentId := r.GetInt("commentId")
	if err := service.Comment.DeleteLabComment(r.Context(), labCommentId); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, true)
}
