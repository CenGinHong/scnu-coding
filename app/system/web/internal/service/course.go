package service

// @Author: 陈健航
// @Date: 2021/1/12 16:39
// @Description:

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/response"
)

var Course = courseService{}

type courseService struct{}

// ListCourseByTeacherId 根据教师id获取该老师所开设的课程信息
// @receiver c *courseService
// @param ctx context.Context
// @param isClose bool
// @return resp *response.PageResp
// @return err error
// @date 2021-08-19 17:24:41
func (c *courseService) ListCourseByTeacherId(ctx context.Context) (resp *response.PageResp, err error) {
	ctxUser := service.Context.Get(ctx).User
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.ListCourseResp, 0)
	d := dao.Course.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(dao.Course.Columns.UserId, ctxUser.UserId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	// 查询
	if err = d.OrderDesc(dao.Course.Columns.CreatedAt).Scan(&records); err != nil {
		return nil, err
	}
	// 分页信息整合
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

// ListCourseEnroll 根据学生id获取该学生修读的课程信息
// @receiver c *courseService
// @param ctx context.Context
// @return resp *response.PageResp
// @return err error
// @date 2021-07-22 19:40:25
func (c *courseService) ListCourseEnroll(ctx context.Context) (resp *response.PageResp, err error) {
	ctxUser := service.Context.Get(ctx).User
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	// 分页
	d := dao.ReCourseUser.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	// 加条件
	d = d.Where(g.Map{
		dao.ReCourseUser.Columns.UserId: ctxUser.UserId,
	})
	// 查总数
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	// 找出所有课程id
	courseIds, err := d.OrderDesc(dao.Course.Columns.CreatedAt).Array(dao.ReCourseUser.Columns.CourseId)
	if err != nil {
		return nil, err
	}
	records := make([]*define.ListCourseResp, 0)
	if err = dao.Course.Ctx(ctx).WherePri(courseIds).WithAll().Scan(&records); err != nil {
		return nil, err
	}
	// 分页信息整合
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

//// UpdateCourse 更新课程
//// @receiver receiver
//// @params ctx
//// @params req
//// @return err
//// @date 2021-05-03 15:52:49
//func (c *courseService) UpdateCourse(ctx context.Context, req *define.UpdateCourseReq) (err error) {
//	ctxUser := service.Context.Get(ctx).User
//	// 保存
//	if _, err = dao.Course.Where(g.Map{
//		dao.Course.Columns.UserId:   ctxUser.UserId,
//		dao.Course.Columns.CourseId: req.CourseId,
//	}).OmitEmpty().Update(req); err != nil {
//		return err
//	}
//	return nil
//}
//
//// InsertCourse 插入课程
//// @receiver receiver
//// @params req
//// @return err
//// @date 2021-05-03 15:52:42
//func (c *courseService) InsertCourse(ctx context.Context, req *define.InsertCourseReq) (err error) {
//	saveModel := gconv.Map(req)
//	ctxUser := service.Context.Get(ctx).User
//	saveModel["UserId"] = ctxUser.UserId
//	// 保存
//	if _, err = dao.Course.InsertLab(saveModel); err != nil {
//		return err
//	}
//	return nil
//}

// SearchCourseByCourseNameOrTeacherName 搜索课程
// @receiver receiver
// @params ctx
// @params courseName
// @return resp
// @return err
// @date 2021-05-03 15:52:29
func (c *courseService) SearchCourseByCourseNameOrTeacherName(ctx context.Context, courseNameOrTeacherName string) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	ctxUserInfo := service.Context.Get(ctx).User
	d := dao.Course.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(dao.Course.Columns.IsClose, false)
	if courseNameOrTeacherName != "" {
		d = d.WhereLike(dao.Course.Columns.CourseName, "%"+courseNameOrTeacherName+"%")
		teacherId, err := dao.SysUser.Ctx(ctx).WhereLike(dao.SysUser.Columns.Username, "%"+courseNameOrTeacherName+"%").
			Array(dao.SysUser.Columns.UserId)
		if err != nil {
			return nil, err
		}
		if len(teacherId) > 0 {
			d = d.WhereOr(dao.Course.Columns.UserId, teacherId)
		}
	}
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	records := make([]*define.SearchCourseResp, 0)
	if err = d.Scan(&records); err != nil {
		return nil, err
	}
	// 找一下有没加入课程
	if err = dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.ReCourseUserDao.Columns.CourseId, gdb.ListItemValuesUnique(records, "CourseId")).
		Where(dao.ReCourseUser.Columns.UserId, ctxUserInfo.UserId).
		Fields("course_id,COUNT(*) as is_take").
		Group("course_id").
		ScanList(&records, "IsTakeDetail", "course_id:CourseId"); err != nil {
		return nil, err
	}
	// 分页信息整合
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

//
//func (c *courseService) Delete(ctx context.Context, courseId int) (err error) {
//	ctxUser := service.Context.Get(ctx).User
//	// 删除课程信息
//	if _, err = dao.Course.Where(g.Map{
//		dao.Course.Columns.UserId:   ctxUser.UserId,
//		dao.Course.Columns.CourseId: courseId,
//	}).Delete(); err != nil {
//		return err
//	}
//	return nil
//}

// ExportCsvTemplate 导出模板
// @receiver receiver
// @return file
// @return err
// @date 2021-05-03 23:59:58
func (c *courseService) ExportCsvTemplate() (file *bytes.Buffer, err error) {
	// 新建csv
	file = &bytes.Buffer{}
	utils.WriteBom(file)
	writer := csv.NewWriter(file)
	defer writer.Flush()
	headLine := make([]string, 0)
	headLine = append(headLine, "学号")
	headLine = append(headLine, "姓名")
	headLine = append(headLine, "班级")
	headLine = append(headLine, "专业")
	if err = writer.Write(headLine); err != nil {
		return nil, err
	}
	return file, nil
}

//func (c *courseService) ImportCsvTemplate(_ context.Context, uploadFile *ghttp.UploadFile, courseId int) (err error) {
//	if gfile.ExtName(uploadFile.Filename) != "csv" {
//		return code.UnSupportUploadTypeError
//	}
//	file, err := uploadFile.Open()
//	if err != nil {
//		return err
//	}
//	defer func(file multipart.File) {
//		_ = file.Close()
//	}(file)
//	// 转格式，去bom头
//	removeBomReader, err := utils.RemoveBom(file)
//	if err != nil {
//		return err
//	}
//	reader := csv.NewReader(removeBomReader)
//	csvData, err := reader.ReadAll()
//	if err != nil {
//		return err
//	}
//	// 开启事务
//	tx, err := g.DB().Begin()
//	if err != nil {
//		return err
//	}
//	// 开启事务
//	defer func() {
//		if err != nil {
//			_ = tx.Rollback()
//		} else {
//			_ = tx.Commit()
//		}
//	}()
//	for i, row := range csvData {
//		// 表头不读
//		if i == 0 {
//			continue
//		}
//		num, realName, class, major := row[0], row[1], row[2], row[3]
//		// 用学校域名邮箱修改
//		email := fmt.Sprintf("%s@%s", num, g.Cfg().GetString("email.domain"))
//		// 找一下该账号是否已存在
//		newStudent := &model.SysUser{}
//		if err = dao.SysUser.TX(tx).WherePri(dao.SysUser.Columns.UserNum, num).Fields(dao.SysUser.Columns.UserId).
//			Scan(&newStudent); err != nil {
//			return err
//		}
//		var stuId int64
//		// 用户不在，帮他注册
//		if newStudent == nil {
//			// 密码加密,用学号当密码
//			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(num), bcrypt.DefaultCost)
//			if err != nil {
//				return err
//			}
//			type UserData struct {
//				Email        string `valid:"required|email#邮箱不能为空|邮箱字段不符合格式"`
//				Num          string `valid:"required|integer#学号不能为空|学号格式不正确"`
//				RealName     string `valid:"required#邮箱不能为空"`
//				Organization string `valid:"required#班级不能为空"`
//				Major        string
//				Password     string `valid:"password#密码不符合强度"`
//			}
//			userData := &UserData{
//				Email:        email,
//				Num:          num,
//				RealName:     realName,
//				Organization: class,
//				Major:        major,
//				Password:     string(hashedPassword),
//			}
//			// 这里有个不知道什么bug，必须用err2,检查字段是否符合
//			if err2 := gvalid.CheckStruct(context.TODO(), userData, nil); err2 != nil {
//				return err2
//			}
//			// 保存
//			stuId, err = dao.SysUser.TX(tx).Data(userData).InsertAndGetId()
//			if err != nil {
//				return err
//			}
//			// 赋予权限
//			if _, err = dao.ReUserRole.TX(tx).InsertLab(g.Map{
//				dao.ReUserRole.Columns.UserId: stuId,
//				dao.ReUserRole.Columns.RoleId: Student,
//			}); err != nil {
//				return err
//			}
//		} else {
//			stuId = int64(newStudent.UserId)
//		}
//		// 加入选课表
//		if _, err = dao.ReCourseUser.TX(tx).Data(g.Map{
//			dao.ReCourseUser.Columns.UserId:   stuId,
//			dao.ReCourseUser.Columns.CourseId: courseId,
//		}).Save(); err != nil {
//			return err
//		}
//		// 还要把签到记录补上
//		checkInRecordIds, err := dao.CheckinRecord.Where(dao.CheckinRecord.Columns.CourseId, courseId).TX(tx).Value(dao.CheckinRecord.Columns.CheckinRecordId)
//		if err != nil {
//			return err
//		}
//		data := make([]g.Map, 0)
//		for _, checkInRecordId := range checkInRecordIds.Ints() {
//			data = append(data, g.Map{
//				dao.CheckinDetail.Columns.CheckinRecordId: checkInRecordId,
//				dao.CheckinDetail.Columns.UserId:          stuId,
//				dao.CheckinDetail.Columns.IsCheckin:       false,
//			})
//		}
//		if _, err = dao.CheckinDetail.Data(data).TX(tx).Save(); err != nil {
//			return err
//		}
//		// 把作业提交补上
//		labIds, err := dao.Lab.Where(dao.Lab.Columns.CourseId, courseId).TX(tx).Value(dao.Lab.Columns.LabId)
//		if err != nil {
//			return err
//		}
//		data = make([]g.Map, 0)
//		for _, labId := range labIds.Ints() {
//			data = append(data, g.Map{
//				dao.LabSubmit.Columns.LabId:  labId,
//				dao.LabSubmit.Columns.UserId: stuId,
//			})
//		}
//		if _, err = dao.LabSubmit.Data(data).TX(tx).Save(); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (c *courseService) GetCourseDetail(ctx context.Context, courseId int) (resp *define.CourseDetailResp, err error) {
	resp = &define.CourseDetailResp{}
	if err = dao.Course.Ctx(ctx).WherePri(courseId).WithAll().Scan(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *courseService) ListCourseStudentOverview(ctx context.Context, courseId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	d := dao.ReCourseUser.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(dao.ReCourseUser.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	records := make([]*define.CourseStudentOverviewResp, 0)
	if err = d.With(define.CourseStudentOverviewResp{}.UserDetail).Scan(&records); err != nil {
		return nil, err
	}
	labIds, err := dao.Lab.Ctx(ctx).Where(dao.Lab.Columns.CourseId, courseId).Array(dao.Lab.Columns.LabId)
	if err != nil {
		return nil, err
	}
	// 查平均成绩
	if err = dao.LabSubmit.Ctx(ctx).Where(dao.LabSubmit.Columns.LabId, labIds).
		Where(dao.LabSubmit.Columns.UserId, gdb.ListItemValuesUnique(records, "UserId")).
		Fields(dao.LabSubmit.Columns.UserId, "AVG(score) as score").Group(dao.LabSubmit.Columns.UserId).
		ScanList(&records, "AvgScoreDetail", "user_id:UserId"); err != nil {
		return nil, err
	}
	// 签到总数
	checkinRecordId, err := dao.CheckinRecord.Ctx(ctx).Where(dao.CheckinRecord.Columns.CourseId, courseId).Array(dao.CheckinRecord.Columns.CheckinRecordId)
	if err != nil {
		return nil, err
	}
	// 签到总数
	checkinCount := len(checkinRecordId)
	for _, record := range records {
		record.CheckinDetail.TotalCount = checkinCount
	}
	if err = dao.CheckinDetail.Ctx(ctx).Where(dao.CheckinDetail.Columns.CheckinRecordId, checkinRecordId).
		Where(dao.CheckinDetail.Columns.UserId, gdb.ListItemValuesUnique(records, "UserId")).
		Where(dao.CheckinDetail.Columns.IsCheckin, true).
		Fields(dao.CheckinDetail.Columns.UserId, "COUNT(*) as checkinCount").
		Group(dao.CheckinDetail.Columns.UserId).
		ScanList(&records, "CheckinDetail", "user_id:UserId"); err != nil {
		return nil, err
	}
	// 查编码时间
	if err = dao.CodingTime.Ctx(ctx).Where(dao.CodingTime.Columns.LabId, labIds).
		Where(dao.CodingTime.Columns.UserId, gdb.ListItemValuesUnique(records, "UserId")).
		Fields(dao.CodingTime.Columns.UserId, fmt.Sprintf("SUM(%s) as codingTime", dao.CodingTime.Columns.Duration)).
		Group(dao.CodingTime.Columns.UserId).
		ScanList(&records, "CodingTimeDetail", "user_id:UserId"); err != nil {
		return nil, err
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

func (c *courseService) ListOneStudentScore(ctx context.Context, courseId int, userId int) (resp *response.PageResp, err error) {
	records := make([]*define.ListOneStudentScore, 0)
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	d := dao.Lab.Ctx(ctx).Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	d = d.Where(dao.Lab.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if err = d.Scan(&records); err != nil {
		return nil, err
	}
	if err = dao.LabSubmit.Ctx(ctx).Where(dao.LabSubmit.Columns.UserId, userId).
		Where(dao.LabSubmit.Columns.LabId, gdb.ListItemValuesUnique(records, "LabId")).
		Fields(define.ListOneStudentScore{}.LabSubmitDetail).
		Scan(&records, "LabSubmitDetail", "lab_id:LabId"); err != nil {
		return nil, err
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

func (c *courseService) IsEnrollCourse(ctx context.Context, courseId int) (isEnroll bool, err error) {
	ctxUser := service.Context.Get(ctx).User
	count, err := dao.ReCourseUser.Ctx(ctx).Where(dao.Course.Columns.CourseId, courseId).Where(dao.Course.Columns.UserId, ctxUser.UserId).Count()
	if err != nil {
		return false, err
	}
	isEnroll = count > 0
	return isEnroll, nil
}

func (c courseService) ImportStudent2Class(ctx context.Context, req *define.ImportStudent2Class) (resp *define.ImportStudent2ClassResp, err error) {
	errorNums := make([]string, 0)
	insertData := make([]g.Map, 0)
	for _, useNum := range req.StudentNums {
		//TODO 学号校验
		count, err := dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns.UserNum, useNum).Count()
		if err != nil {
			return nil, err
		}
		//该用户未注册
		if count == 0 {
			errorNums = append(errorNums, useNum)
			continue
		}
		// 查id
		userId, err := dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns.UserNum, useNum).Value(dao.SysUser.Columns.UserId)
		if err != nil {
			return nil, err
		}
		// 加入插入数据
		insertData = append(insertData, g.Map{
			dao.ReCourseUser.Columns.UserId:   userId,
			dao.ReCourseUser.Columns.CourseId: req.CourseId,
		})
	}
	confirmStudent := make([]*define.ConfirmStudentResp, 0)
	if _, err = dao.ReCourseUser.ReCourseUserDao.Ctx(ctx).Data(insertData).Batch(len(insertData)).Save(); err != nil {
		return nil, err
	}
	if err = dao.SysUser.Ctx(ctx).WherePri(gdb.ListItemValuesUnique(insertData, dao.ReCourseUser.Columns.UserId)).
		Scan(&confirmStudent); err != nil {
		return nil, err
	}
	resp = &define.ImportStudent2ClassResp{ErrorStudentNums: errorNums, SuccessRecords: confirmStudent}
	return resp, err
}

func (c *courseService) IsOpenByTeacherId(ctx context.Context, courseId int) (isOpen bool, err error) {
	ctxUser := service.Context.Get(ctx).User
	count, err := dao.Course.Ctx(ctx).WherePri(courseId).Where(dao.Course.Columns.UserId, ctxUser.UserId).Count()
	if err != nil {
		return false, err
	}
	isOpen = count > 0
	return isOpen, nil
}

func (c *courseService) RemoveStudentFromClass() {

}
