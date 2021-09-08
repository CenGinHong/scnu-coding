package service

// @Author: 陈健航
// @Date: 2021/1/15 21:30
// @Description:

import (
	"context"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/library/response"

	"github.com/gogf/gf/frame/g"
)

var Comment = commentService{}

type commentService struct{}

// ListCourseCommentByCourseId 列出所有课程评论
// @receiver s
// @params ctx
// @params courseId
// @return resp
// @return error
// @date 2021-05-04 23:58:23
func (c *commentService) ListCourseCommentByCourseId(ctx context.Context, courseId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.CourseCommentResp, 0)
	// 查询主评
	d := dao.CourseComment.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(g.Map{dao.CourseComment.Columns.CourseId: courseId, dao.CourseComment.Columns.Pid: 0})
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if err = d.OrderDesc(dao.CourseComment.Columns.CreatedAt).WithAll().Scan(&records); err != nil {
		return nil, err
	}
	// 分页信息整合
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

// ListLabCommentByLabId 列出所有实验评论
// @receiver s
// @params ctx
// @params labId
// @return resp
// @return error
// @date 2021-05-04 23:58:17
func (c *commentService) ListLabCommentByLabId(ctx context.Context, labId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.LabCommentResp, 0)
	d := dao.LabComment.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(g.Map{dao.LabComment.Columns.LabId: labId, dao.LabComment.Columns.Pid: 0})
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if err = d.OrderDesc(dao.LabComment.Columns.CreatedAt).WithAll().Scan(&records); err != nil {
		return nil, err
	}
	// 分页信息整合
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

// InsertCourseComment 插入课程评论
// @receiver s
// @params ctx
// @params req
// @return err
// @date 2021-05-04 23:58:11
func (c *commentService) InsertCourseComment(ctx context.Context, req *define.InsertCourseCommentReq) (id int64, err error) {
	ctxUser := service.Context.Get(ctx).User
	// 保存模型
	courseComment := g.Map{
		dao.CourseComment.Columns.CourseId:    req.CourseId,
		dao.CourseComment.Columns.CommentText: req.CommentText,
		dao.CourseComment.Columns.Pid:         req.Pid,
		dao.CourseComment.Columns.UserId:      ctxUser.UserId,
	}
	// 保存
	if id, err = dao.CourseComment.Ctx(ctx).Data(courseComment).InsertAndGetId(); err != nil {
		return 0, err
	}
	return id, nil
}

// InsertLabComment 插入实验评论
// @receiver s
// @params ctx
// @params req
// @return error
// @date 2021-05-04 23:57:52
func (c *commentService) InsertLabComment(ctx context.Context, req *define.InsertLabCommentReq) (id int64, err error) {
	ctxUser := service.Context.Get(ctx).User
	labComment := g.Map{
		dao.LabComment.Columns.CommentText: req.CommentText,
		dao.LabComment.Columns.UserId:      ctxUser.UserId,
		dao.LabComment.Columns.LabId:       req.LabId,
		dao.LabComment.Columns.Pid:         req.Pid,
	}
	// 保存
	if id, err = dao.LabComment.Ctx(ctx).Data(labComment).InsertAndGetId(); err != nil {
		return 0, err
	}
	return id, nil
}

// DeleteLabComment 删除实验评论
// @receiver s
// @params ctx
// @params commentId
// @return err
// @date 2021-05-04 23:57:45
func (c *commentService) DeleteLabComment(ctx context.Context, commentId int) (err error) {
	delData := make([]int, 0)
	delData = append(delData, commentId)
	newDelData, err := dao.LabComment.Ctx(ctx).Where(dao.LabComment.Columns.Pid, commentId).Array(dao.LabComment.Columns.CommentId)
	if err != nil {
		return err
	}
	for len(newDelData) > 0 {
		// 追加新的删除数据
		delData = append(delData, gconv.Ints(newDelData)...)
		newDelData, err = dao.LabComment.Ctx(ctx).Where(dao.LabComment.Columns.Pid, newDelData).Array(dao.LabComment.Columns.CommentId)
		if err != nil {
			return err
		}
	}
	if _, err = dao.LabComment.Ctx(ctx).WherePri(delData).Delete(); err != nil {
		return err
	}
	return nil
}

// DeleteCourseComment 删除课程评论
// @receiver s
// @params ctx
// @params commentId
// @return err
// @date 2021-05-04 23:57:40
func (c *commentService) DeleteCourseComment(ctx context.Context, commentId int) (err error) {
	delData := make([]int, 0)
	delData = append(delData, commentId)
	newDelData, err := dao.CourseComment.Ctx(ctx).Where(dao.CourseComment.Columns.Pid, commentId).Array(dao.CourseComment.Columns.CommentId)
	if err != nil {
		return err
	}
	for len(newDelData) > 0 {
		// 追加新的删除数据
		delData = append(delData, gconv.Ints(newDelData)...)
		newDelData, err = dao.CourseComment.Ctx(ctx).Where(dao.CourseComment.Columns.Pid, newDelData).Array(dao.CourseComment.Columns.CommentId)
		if err != nil {
			return err
		}
	}
	if _, err = dao.CourseComment.Ctx(ctx).WherePri(delData).Delete(); err != nil {
		return err
	}
	return nil
}
