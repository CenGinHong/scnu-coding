package service

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/library/response"
	"scnu-coding/library/str"
)

// @Author: 陈健航
// @Date: 2021/2/10 22:34
// @Description:

var Lab = labService{}

type labService struct{}

// InsertLab 新建实验
// @receiver s
// @params _
// @params req
// @return err
// @return err
// @date 2021-05-05 00:14:12
func (l *labService) InsertLab(ctx context.Context, req *define.InsertLabReq) (labId int64, err error) {
	labId, err = dao.Lab.Ctx(ctx).Data(req).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return labId, nil
}

// ListLabByCourseId 课程详情页查询分页列表实验
// @receiver l *labService
// @param ctx context.Context
// @param courseId int
// @return resp *response.PageResp
// @return err error
// @date 2021-08-02 22:01:08
func (l *labService) ListLabByCourseId(ctx context.Context, courseId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	d := dao.Lab.Ctx(ctx)
	if ctxPageInfo != nil {
		d = d.Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	}
	d = d.Where(dao.Lab.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	d = d.OrderAsc(dao.Lab.Columns.CreatedAt)
	records := make([]*define.LabDetailResp, 0)
	if err = d.With(define.LabDetailResp{}.LabSubmitDetail).Scan(&records); err != nil {
		return nil, err
	}
	// 拼接地址
	addr := service.File.GetMinioAddr(ctx)
	for _, record := range records {
		if record.AttachmentSrc != "" {
			record.AttachmentSrc = addr + "/" + record.AttachmentSrc
		}
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

// Update 更新实验
// @receiver s
// @params _
// @params req
// @return err
// @date 2021-05-05 13:11:21
func (l *labService) Update(ctx context.Context, req *define.UpdateLabReq) (err error) {
	if _, err = dao.Lab.Ctx(ctx).OmitNilData().WherePri(req.LabId).Update(req); err != nil {
		return err
	}
	return nil
}

// Delete 删除实验
// @receiver s
// @params ctx
// @params labId
// @return err
// @date 2021-05-05 13:11:31
func (l *labService) Delete(ctx context.Context, labId int) (err error) {
	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if _, err = dao.Lab.Ctx(ctx).TX(tx).WherePri(labId).Delete(); err != nil {
			return err
		}
		// 删除所有提交的实验
		if _, err = dao.LabSubmit.Ctx(ctx).TX(tx).Where(dao.LabSubmit.Columns.LabId, labId).Delete(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (l *labService) isTake(ctx context.Context, labId int) (isTake bool, err error) {
	ctxUser := service.Context.Get(ctx).User
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Cache(0).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return false, err
	}
	count, err := dao.ReCourseUser.Ctx(ctx).Cache(0, str.CACHE_ENROLL).Where(dao.ReCourseUser.Columns.CourseId, courseId).
		Where(dao.ReCourseUser.Columns.UserId, ctxUser.UserId).Count()
	if err != nil {
		return false, err
	}
	isTake = count > 0
	return isTake, nil
}

func (l *labService) teacherIsTake(ctx context.Context, labId int) (isTake bool, err error) {
	ctxUser := service.Context.Get(ctx).User
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Cache(0).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return false, err
	}
	count, err := dao.Course.Ctx(ctx).WherePri(courseId).Where(dao.Course.Columns.UserId, ctxUser.UserId).Count()
	if err != nil {
		return false, err
	}
	isTake = count > 0
	return isTake, nil
}
