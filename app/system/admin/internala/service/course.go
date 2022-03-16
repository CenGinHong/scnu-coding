package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/admin/internala/define"
	"scnu-coding/library/response"
	"strings"
)

var Course = courseService{}

type courseService struct {
}

func (s *courseService) ListAllCourse(ctx context.Context) (resp *response.PageResp, err error) {
	// 获取分页信息
	pageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.ListCourseResp, 0)
	// 筛选项
	d := dao.Course.Ctx(ctx)
	for k, v := range pageInfo.ParseFilterFields {
		d = d.Where(k, v)
	}
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if pageInfo != nil {
		d = d.Page(pageInfo.Current, pageInfo.PageSize)
	}
	if err = d.WithAll().Scan(&records); err != nil {
		return nil, err
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

func (s *courseService) ListEnroll(ctx context.Context, courseId int) (resp *response.PageResp, err error) {
	// 获取分页信息
	pageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.CourseEnroll, 0)
	// 所有的学生学号
	stuIds, err := dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.Columns.CourseId, courseId).Array(dao.ReCourseUser.Columns.UserId)
	if err != nil {
		return nil, err
	}
	d := dao.SysUser.Ctx(ctx)
	d = d.WherePri(stuIds)
	if pageInfo != nil {
		// 排序
		if pageInfo.SortField != "" {
			d = d.Order(pageInfo.SortField, pageInfo.SortOrder)
		}
		// 加入筛选条件
		for key, value := range pageInfo.ParseFilterFields {
			d = d.Where(key, value)
		}
	}
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if err = d.Scan(&records); err != nil {
		return nil, err
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

func (s *courseService) RemoveCourseEnroll(ctx context.Context, req *define.RemoveCourseEnrollReq) (err error) {
	if _, err = dao.ReCourseUser.Ctx(ctx).Where(g.Map{
		dao.ReCourseUser.Columns.UserId:   req.UserIds,
		dao.ReCourseUser.Columns.CourseId: req.CourseId,
	}).Delete(); err != nil {
		return err
	}
	return nil
}

func (s *courseService) AddStudent2Class(ctx context.Context, req *define.AddStudent2ClassReq) (errMsg string, err error) {
	data := make([]g.Map, 0)
	sb := &strings.Builder{}
	for _, r := range req.StudentNums {
		// 找到对应的id
		userId, err := dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns.UserNum, r).Value(dao.SysUser.Columns.UserId)
		if err != nil {
			return "", err
		}
		if userId.IsNil() {
			sb.WriteString(fmt.Sprintf("账号不存在：%s\n", r))
		} else {
			data = append(data, g.Map{
				dao.ReCourseUser.Columns.CourseId: req.CourseId,
				dao.ReCourseUser.Columns.UserId:   userId,
			})
		}
	}
	if len(data) > 0 {
		if _, err = dao.ReCourseUser.Ctx(ctx).Batch(len(data)).Data(data).InsertIgnore(); err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

//func (s *courseService) ListAllCourse(ctx context.Context) (resp *response.PageResp, err error) {
//	// 获取分页信息
//	pageInfo := service.Context.Get(ctx).PageInfo
//	records := make([]*define.ListCourseResp, 0)
//	// 筛选集
//	filter := make(map[string][]*response.FilterType, 0)
//	// 查找可筛选项 教师
//	roleId, err := dao.SysRole.Ctx(ctx).WhereNot(dao.SysRole.Columns.Description, "Student").FindArray(dao.SysRole.Columns.RoleId)
//	if err != nil {
//		return nil, err
//	}
//	teacherDetail := make([]*model.SysUser, 0)
//	if err = dao.SysUser.Ctx(ctx).Distinct().Fields(dao.SysUser.Columns.UserId, dao.SysUser.Columns.Username).
//		Where(dao.SysUser.Columns.RoleId, roleId).Scan(&teacherDetail); err != nil {
//		return nil, err
//	}
//	tempFilter := make([]*response.FilterType, 0)
//	for _, value := range teacherDetail {
//		tempFilter = append(tempFilter, &response.FilterType{
//			Text:     value.Username,
//			Value:    gconv.String(value.UserId),
//			Children: nil,
//		})
//	}
//	filter["teacher"] = tempFilter
//	// 查找可筛选项 语言
//	LanguageTypeValue, err := dao.Course.Ctx(ctx).Distinct().FindArray(dao.Course.Columns.LanguageType)
//	if err != nil {
//		return nil, err
//	}
//	tempFilter = make([]*response.FilterType, 0)
//	for _, value := range LanguageTypeValue {
//		tempFilter = append(tempFilter, &response.FilterType{
//			Text:     language_enum.Num2LanguageString(value.Int()),
//			Value:    value.String(),
//			Children: nil,
//		})
//	}
//	// 可筛选项 是否结课
//	filter["languageType"] = tempFilter
//	isCloseValue, err := dao.Course.Ctx(ctx).Distinct().FindArray(dao.Course.Columns.IsClose)
//	if err != nil {
//		return nil, err
//	}
//	tempFilter = make([]*response.FilterType, 0)
//	for _, value := range isCloseValue {
//		tempFilter = append(tempFilter, &response.FilterType{
//			Text:     enum.Num2IsCloseString(value.Int()),
//			Value:    value.String(),
//			Children: nil,
//		})
//	}
//	filter["isClose"] = tempFilter
//	// 筛选项
//	d := dao.Course.Ctx(ctx).Page(pageInfo.Current, pageInfo.PageSize)
//	for k, v := range pageInfo.ParseFilterFields {
//		d = d.Where(k, v)
//	}
//	//查总数
//	total, err := d.Count()
//	if err != nil {
//		return nil, err
//	}
//	// 排序项
//	d = d.Order(pageInfo.SortField, pageInfo.SortOrder)
//	// 查询
//	if err = d.Fields(define.Course{}).WithAll().Scan(&records); err != nil {
//		return nil, err
//	}
//	// 统计选课人数
//	for _, record := range records {
//		enrollCount, err := dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.Columns.CourseId, record.CourseId).Count()
//		if err != nil {
//			return nil, err
//		}
//		record.EnrollCount = enrollCount
//	}
//	resp = response.GetPageResp(records, total, filter)
//	return resp, nil
//}
