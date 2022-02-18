package service

// @Author: 陈健航
// @Date: 2021/1/16 22:04
// @Description:

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/response"
	"time"
)

var Checkin = checkinService{
	checkinKeyCache: utils.NewMyCache(),
}

type checkinService struct {
	checkinKeyCache *utils.MyCache
}

// ListCheckinRecordByCourseId 教师获取签到列表
// @receiver c *checkinService
// @param ctx context.Context
// @param courseId int
// @return resp *response.PageResp
// @return err error
// @date 2021-08-09 23:15:50
func (c *checkinService) ListCheckinRecordByCourseId(ctx context.Context, courseId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	d := dao.CheckinRecord.Ctx(ctx)
	d = d.Where(dao.CheckinRecord.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	records := make([]*define.CheckinRecordResp, 0)
	if ctxPageInfo.SortOrder != "" {
		d = d.Order(ctxPageInfo.SortField, ctxPageInfo.SortOrder)
	}
	if err = d.Page(ctxPageInfo.Current, ctxPageInfo.PageSize).Scan(&records); err != nil {
		return nil, err
	}
	// 课程选课的人数
	totalTakeCount, err := dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.Columns.CourseId, courseId).Count()
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		record.Attendance.TakeCount = totalTakeCount
	}
	// 实际参与签到的人
	if err = dao.CheckinDetail.Ctx(ctx).Where(g.Map{
		dao.CheckinDetail.Columns.CheckinRecordId: gdb.ListItemValuesUnique(records, "CheckinRecordId"),
		dao.CheckinDetail.Columns.IsCheckin:       true,
	}).Fields(dao.CheckinDetail.Columns.CheckinRecordId, "COUNT(*) as checkin_count").
		Group(dao.CheckinDetail.Columns.CheckinRecordId).
		ScanList(&records, "Attendance", "checkin_record_id:CheckinRecordId"); err != nil {
		return nil, err
	}
	return response.GetPageResp(records, total, nil), nil
}

// StartCheckin 教师发起签到
// @receiver c *checkinService
// @param _ context.Context
// @param req *define.StartCheckInReq
// @return err error
// @date 2021-08-03 20:29:42
func (c *checkinService) StartCheckin(ctx context.Context, req *define.StartCheckInReq) (checkinRecordId int64, err error) {
	checkinRecordId, err = dao.CheckinRecord.Ctx(ctx).Data(g.Map{
		dao.CheckinRecord.Columns.CheckinKey:  req.CheckinKey,
		dao.CheckinRecord.Columns.CheckinName: req.CheckinName,
		dao.CheckinRecord.Columns.CourseId:    req.CourseId,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	// 存入签到密钥,限时
	cacheData := &define.RedisCheckinData{
		CheckinName:     req.CheckinName,
		CheckinKey:      req.CheckinKey,
		CheckinRecordId: checkinRecordId,
		TotalDuration:   req.Duration,
	}
	if err = c.checkinKeyCache.Set(req.CourseId, cacheData, time.Duration(req.Duration)*time.Second); err != nil {
		return 0, err
	}
	return checkinRecordId, nil
}

// GetCheckinStatus 获取签到进行时状态
// @receiver c *checkinService
// @param _ context.Context
// @param courseId int
// @return resp *define.CheckinStatusResp
// @return err error
// @date 2021-08-03 20:30:21
func (c *checkinService) GetCheckinStatus(_ context.Context, courseId int) (resp *define.CheckinStatusResp, err error) {
	v, err := c.checkinKeyCache.GetVar(courseId)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		return nil, nil
	}
	resp = &define.CheckinStatusResp{}
	if err = v.Struct(&resp); err != nil {
		return nil, err
	}
	// 获取未过期的时间
	expire, err := c.checkinKeyCache.GetExpire(courseId)
	if err != nil {
		return nil, err
	}
	resp.RemainDuration = expire.Seconds()
	return resp, nil
}

// StuListCheckinRecords 学生列表自己签到记录
// @receiver c *checkinService
// @param ctx context.Context
// @param courseId int
// @return resp *response.PageResp
// @return err error
// @date 2021-08-04 08:49:30
func (c *checkinService) StuListCheckinRecords(ctx context.Context, courseId int) (resp *response.PageResp, err error) {
	records := make([]*define.StuListCheckInRecordResp, 0)
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	ctxUser := service.Context.Get(ctx).User
	d := dao.CheckinRecord.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(dao.CheckinRecord.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if ctxPageInfo.SortOrder != "" {
		d = d.Order(ctxPageInfo.SortField, ctxPageInfo.SortOrder)
	}
	if err = d.Scan(&records); err != nil {
		return nil, err
	}
	// 绑定详细信息
	if err = dao.CheckinDetail.Ctx(ctx).
		Where(dao.CheckinDetail.Columns.CheckinRecordId, gdb.ListItemValuesUnique(records, "CheckinRecordId")).
		Where(dao.CheckinDetail.Columns.UserId, ctxUser.UserId).Fields(define.StuListCheckInRecordResp{}.CheckinDetail).
		ScanList(records, "CheckInDetail", "checkin_record_id:CheckinRecordId"); err != nil {
		return nil, err
	}
	// 分页信息整合
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

// ListCheckinDetailByCheckInRecordId 教师列表签到详情记录
// @receiver receiver
// @params ctx
// @params checkInRecordId
// @return resp
// @return err
// @date 2021-05-05 23:41:06
func (c *checkinService) ListCheckinDetailByCheckInRecordId(ctx context.Context, checkInRecordId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.CheckinDetailResp, 0)
	//找出所有加入课程的学生
	courseId, err := dao.CheckinRecord.Ctx(ctx).WherePri(checkInRecordId).Value(dao.CheckinRecord.Columns.CourseId)
	if err != nil {
		return nil, err
	}
	d := dao.ReCourseUser.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(dao.ReCourseUser.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	//查出选课学生的个人资料
	if err = d.With(define.CheckinDetailResp{}.UserDetail).Scan(&records); err != nil {
		return nil, err
	}
	// c查出有无参与该签到
	if err = dao.CheckinDetail.Ctx(ctx).Where(dao.CheckinDetail.Columns.CheckinRecordId, checkInRecordId).
		Where(dao.CheckinDetail.Columns.UserId, gdb.ListItemValuesUnique(records, "UserId")).
		ScanList(&records, "CheckinDetail", "user_id:UserId"); err != nil {
		return nil, err
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

func (c *checkinService) UpdateCheckinDetail(ctx context.Context, req *define.UpdateCheckinDetailReq) (err error) {
	if _, err = dao.CheckinDetail.Ctx(ctx).Data(g.Map{
		dao.CheckinDetail.Columns.CheckinRecordId: req.CheckinRecordId,
		dao.CheckinDetail.Columns.IsCheckin:       req.IsCheckin,
		dao.CheckinDetail.Columns.UserId:          req.UserId,
	}).Save(); err != nil {
		return err
	}
	return nil
}

// CheckIn 学生完成签到
// @receiver c
// @params req
// @return err
// @date 2021-03-16 14:24:58
func (c *checkinService) CheckIn(ctx context.Context, req *define.StudentCheckinReq) (err error) {
	ctxUser := service.Context.Get(ctx).User
	// 获取签到密钥
	v, err := c.checkinKeyCache.GetVar(req.CourseId)
	if err != nil {
		return err
	}
	if v.IsNil() {
		return gerror.NewCode(gcode.CodeOperationFailed, "签到已结束")
	}
	cacheData := &define.RedisCheckinData{}
	if err = v.Struct(&cacheData); err != nil {
		return err
	}
	if req.CheckinKey != cacheData.CheckinKey {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "签到密钥错误")
	}
	// 签到码正确,写入数据库
	if _, err = dao.CheckinDetail.Ctx(ctx).Where(g.Map{
		dao.CheckinDetail.Columns.UserId:          ctxUser.UserId,
		dao.CheckinDetail.Columns.CheckinRecordId: cacheData.CheckinRecordId,
	}).Data(dao.CheckinDetail.Columns.IsCheckin, true).Update(); err != nil {
		return err
	}
	return nil
}

func (c *checkinService) DeleteCheckinRecord(ctx context.Context, id int) (err error) {
	if err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		if _, err = dao.CheckinRecord.Ctx(ctx).TX(tx).WherePri(id).Delete(); err != nil {
			return err
		}
		if _, err = dao.CheckinDetail.Ctx(ctx).TX(tx).Where(dao.CheckinDetail.Columns.CheckinRecordId, id).Delete(); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *checkinService) ExportCheckinCsv(ctx context.Context, courseId int) (file *bytes.Buffer, err error) {
	userRecords := make([]*define.EnrollUserDetail, 0)
	// 查出全部选课学生
	if err = dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.Columns.CourseId, courseId).WithAll().
		Scan(&userRecords); err != nil {
		return nil, err
	}
	checkinRecords := make([]*define.ExportCheckinRecord, 0)
	if err = dao.CheckinRecord.Ctx(ctx).Where(dao.CheckinRecord.Columns.CourseId, courseId).WithAll().
		Scan(&checkinRecords); err != nil {
		return nil, err
	}
	// 新建csv
	file = &bytes.Buffer{}
	utils.WriteBom(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headLine := []string{"姓名", "学号"}
	// 写入表头
	for _, record := range checkinRecords {
		headLine = append(headLine, record.CheckinName)
	}
	if err = writer.Write(headLine); err != nil {
		return nil, err
	}
	data := make([][]string, 0)
	// 数据量也不大就直接三层for了
	for _, record := range userRecords {
		row := make([]string, 0)
		row = append(row, record.UserDetail.Username)
		row = append(row, record.UserDetail.UserNum)
		for _, checkinRecord := range checkinRecords {
			isCheckin := false
			count := 0
			for _, checkinDetail := range checkinRecord.CheckinDetails {
				if checkinDetail.UserId == record.UserId && checkinDetail.IsCheckin {
					count++
					isCheckin = true
					row = append(row, "√")
					break
				}
			}
			if !isCheckin {
				row = append(row, " ")
			}
		}
		data = append(data, row)
	}
	if err = writer.WriteAll(data); err != nil {
		return nil, err
	}
	return file, nil
}
