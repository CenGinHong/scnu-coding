package service

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/admin/internala/define"
	"scnu-coding/library/response"
)

var LabSubmit = labSubmitService{}

type labSubmitService struct {
}

func (s labSubmitService) ListLabSubmitByLabId(ctx context.Context, labId int) (resp *response.PageResp, err error) {
	// 获取分页信息
	pageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.ListLabSubmitByLabIdResp, 0)
	// 根据实验id找课程
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return nil, err
	}
	// 找出所有的学生id
	d := dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.Columns.CourseId, courseId)
	if pageInfo != nil {
		d = d.Page(pageInfo.Current, pageInfo.PageSize)
	}
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if err = d.WithAll().Scan(&records); err != nil {
		return nil, err
	}
	if err = dao.LabSubmit.Ctx(ctx).Where(dao.LabSubmit.Columns.LabId, labId).Where(dao.LabSubmit.Columns.UserId,
		gdb.ListItemValuesUnique(records, "UserId")).Fields(define.ListLabSubmitByLabIdResp{}.LabSubmitDetail).
		ScanList(&records, "LabSubmitDetail", "user_id:UserId"); err != nil {
		return nil, err
	}

	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}
