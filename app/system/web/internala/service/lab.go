package service

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/library/response"
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
func (l *labService) InsertLab(ctx context.Context, req *define.InsertLabReq) (err error) {
	// 如果文件就插入文件
	if req.UploadFile != nil {
		req.AttachmentSrc, err = service.File.UploadFile(ctx, req.UploadFile)
		if err != nil {
			return err
		}
		defer func() {
			if err != nil {
				_ = service.File.RemoveObject(ctx, req.AttachmentSrc)
			}
		}()
	}
	// 插入新数据
	if _, err = dao.Lab.Ctx(ctx).
		Data(req).
		Insert(); err != nil {
		return err
	}
	return nil

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
	for _, r := range records {
		if r.AttachmentSrc != "" {
			r.AttachmentSrc = service.File.GetObjectUrl(ctx, r.AttachmentSrc)
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
	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		// 移除旧的文件
		if req.IsRemoveFile {
			attachmentSrc, err := dao.Lab.Ctx(ctx).TX(tx).
				WherePri(req.LabId).
				Value(dao.Lab.Columns.AttachmentSrc)
			if err != nil {
				return err
			}
			defer func() {
				if err != nil {
					_ = service.File.RemoveObject(ctx, attachmentSrc.String())
				}
			}()
			if _, err = dao.Lab.Ctx(ctx).TX(tx).
				WherePri(req.LabId).
				Data(g.Map{dao.Lab.Columns.AttachmentSrc: ""}).
				Update(); err != nil {
				return err
			}
		}
		if req.UploadFile != nil {
			if req.AttachmentSrc, err = service.File.UploadFile(ctx, req.UploadFile); err != nil {
				return err
			}
		}
		if _, err = dao.Lab.Ctx(ctx).TX(tx).
			WherePri(req.LabId).
			Update(req); err != nil {
			return err
		}
		return nil
	}); err != nil {
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
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return false, err
	}
	count, err := dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.Columns.CourseId, courseId).
		Where(dao.ReCourseUser.Columns.UserId, ctxUser.UserId).Count()
	if err != nil {
		return false, err
	}
	isTake = count > 0
	return isTake, nil
}

func (l *labService) teacherIsTake(ctx context.Context, labId int) (isTake bool, err error) {
	ctxUser := service.Context.Get(ctx).User
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
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

func (l *labService) GetOne(ctx context.Context, labId int) (resp *define.LabDetailResp, err error) {
	if err = dao.Lab.Ctx(ctx).WherePri(labId).Scan(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}
