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
// @Date: 2021/3/1 0:36
// @Description:

var CourseAnnouncement = courseAnnouncementService{}

type courseAnnouncementService struct{}

// InsertCourseAnnouncement 插入课程公告
// @receiver receiver
// @params _
// @params req
// @return err
// @date 2021-05-08 09:25:57
func (c *courseAnnouncementService) InsertCourseAnnouncement(ctx context.Context, req *define.InsertCourseAnnouncementReq) (err error) {
	if req.UploadFile != nil {
		if req.AttachmentSrc, err = service.File.UploadFile(ctx, req.UploadFile); err != nil {
			return err
		}
		defer func() {
			if err != nil {
				_ = service.File.RemoveObject(ctx, req.AttachmentSrc)
			}
		}()
	}
	// 插入新数据
	if _, err = dao.CourseAnnouncement.Ctx(ctx).Insert(req); err != nil {
		return err
	}
	return nil
}

// UpdateCourseAnnouncement 更新公告
// @receiver receiver
// @params ctx
// @params req
// @return err
// @date 2021-05-08 09:26:04
func (c *courseAnnouncementService) UpdateCourseAnnouncement(ctx context.Context, req *define.UpdateCourseAnnouncementReq) (err error) {
	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		// 移除旧的文件
		if req.IsRemoveFile {
			attachmentSrc, err := dao.CourseAnnouncement.Ctx(ctx).TX(tx).
				WherePri(req.CourseAnnouncementId).
				Value(dao.Lab.Columns.AttachmentSrc)
			if err != nil {
				return err
			}
			defer func() {
				if err != nil {
					_ = service.File.RemoveObject(ctx, attachmentSrc.String())
				}
			}()
			if _, err = dao.CourseAnnouncement.Ctx(ctx).TX(tx).
				WherePri(req.CourseAnnouncementId).
				Data(g.Map{dao.Lab.Columns.AttachmentSrc: ""}).
				Update(); err != nil {
				return err
			}
		}
		if req.UploadFile != nil {
			req.AttachmentSrc, err = service.File.UploadFile(ctx, req.UploadFile)
			if err != nil {
				return err
			}
		}
		if _, err = dao.CourseAnnouncement.Ctx(ctx).TX(tx).
			WherePri(req.CourseAnnouncementId).
			Update(req); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

// ListCourseAnnouncement 列表课程公告
// @receiver receiver
// @params ctx
// @params courseId
// @return resp
// @return err
// @date 2021-05-08 09:26:18
func (c *courseAnnouncementService) ListCourseAnnouncement(ctx context.Context, courseId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	d := dao.CourseAnnouncement.Ctx(ctx)
	d.Where(dao.CourseAnnouncement.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	records := make([]*define.CourseAnnouncementResp, 0)
	if err = d.Page(ctxPageInfo.Current, ctxPageInfo.PageSize).
		OrderDesc(dao.CourseAnnouncement.Columns.CreatedAt).
		WithAll().
		Scan(&records); err != nil {
		return nil, err
	}
	for _, r := range records {
		r.AttachmentSrc = service.File.GetObjectUrl(ctx, r.AttachmentSrc)
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

// DeleteCourseResource 删除课程公告
// @receiver receiver
// @params ctx
// @params courseRecourseId
// @return err
// @date 2021-05-08 09:26:25
func (c *courseAnnouncementService) DeleteCourseResource(ctx context.Context, courseRecourseId int) (err error) {
	if _, err = dao.CourseAnnouncement.Ctx(ctx).WherePri(courseRecourseId).Delete(); err != nil {
		return err
	}
	return nil
}

func (c *courseAnnouncementService) GetOne(ctx context.Context, courseAnnouncementId int) (resp *define.CourseAnnouncementResp, err error) {
	if err = dao.CourseAnnouncement.Ctx(ctx).WherePri(courseAnnouncementId).Scan(&resp); err != nil {
		return nil, err
	}
	resp.AttachmentSrc = service.File.GetObjectUrl(ctx, resp.AttachmentSrc)
	return resp, nil
}

//// GetOne 获得单个资源
//// @receiver receiver
//// @params _
//// @params courseRecourseId
//// @return resp
//// @return err
//// @date 2021-05-08 09:26:45
//func (c *courseAnnouncementService) GetOne(_ context.Context, courseRecourseId int) (resp *define.CourseAnnouncementResp, err error) {
//	resp = new(define.CourseAnnouncementResp)
//	if err = dao.CourseAnnouncement.WherePri(courseRecourseId).WithAll().Scan(&resp); err != nil {
//		return nil, err
//	}
//	return resp, nil
//}
