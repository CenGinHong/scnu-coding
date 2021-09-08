package service

import (
	"context"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/model"
	"scnu-coding/app/service"
	"scnu-coding/app/system/admin/internal/define"
	"scnu-coding/library/enum"
	"scnu-coding/library/response"
)

var Course = courseService{}

type courseService struct {
}

func (s *courseService) GetAllCourse(ctx context.Context) (resp *response.PageResp, err error) {
	// 获取分页信息
	pageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.Course, 0)
	// 筛选集
	filter := make(map[string][]*response.FilterType, 0)
	// 查找可筛选项 教师
	roleId, err := dao.SysRole.Ctx(ctx).WhereNot(dao.SysRole.Columns.Description, "Student").FindArray(dao.SysRole.Columns.RoleId)
	if err != nil {
		return nil, err
	}
	teacherDetail := make([]*model.SysUser, 0)
	if err = dao.SysUser.Ctx(ctx).Distinct().Fields(dao.SysUser.Columns.UserId, dao.SysUser.Columns.Username).
		Where(dao.SysUser.Columns.RoleId, roleId).Scan(&teacherDetail); err != nil {
		return nil, err
	}
	tempFilter := make([]*response.FilterType, 0)
	for _, value := range teacherDetail {
		tempFilter = append(tempFilter, &response.FilterType{
			Text:     value.Username,
			Value:    gconv.String(value.UserId),
			Children: nil,
		})
	}
	filter["teacher"] = tempFilter
	// 查找可筛选项 语言
	LanguageTypeValue, err := dao.Course.Ctx(ctx).Distinct().FindArray(dao.Course.Columns.LanguageType)
	if err != nil {
		return nil, err
	}
	tempFilter = make([]*response.FilterType, 0)
	for _, value := range LanguageTypeValue {
		tempFilter = append(tempFilter, &response.FilterType{
			Text:     enum.Num2LanguageString(value.Int()),
			Value:    value.String(),
			Children: nil,
		})
	}
	// 可筛选项 是否结课
	filter["languageType"] = tempFilter
	isCloseValue, err := dao.Course.Ctx(ctx).Distinct().FindArray(dao.Course.Columns.IsClose)
	if err != nil {
		return nil, err
	}
	tempFilter = make([]*response.FilterType, 0)
	for _, value := range isCloseValue {
		tempFilter = append(tempFilter, &response.FilterType{
			Text:     enum.Num2IsCloseString(value.Int()),
			Value:    value.String(),
			Children: nil,
		})
	}
	filter["isClose"] = tempFilter
	// 筛选项
	d := dao.Course.Ctx(ctx).Page(pageInfo.Current, pageInfo.PageSize)
	for k, v := range pageInfo.ParseFilterFields {
		d = d.Where(k, v)
	}
	//查总数
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	// 排序项
	d = d.Order(pageInfo.SortField, pageInfo.SortOrder)
	// 查询
	if err = d.Fields(define.Course{}).WithAll().Scan(&records); err != nil {
		return nil, err
	}
	// 统计选课人数
	for _, record := range records {
		enrollCount, err := dao.ReCourseUser.Ctx(ctx).Where(dao.ReCourseUser.Columns.CourseId, record.CourseId).Count()
		if err != nil {
			return nil, err
		}
		record.EnrollCount = enrollCount
	}
	resp = response.GetPageResp(records, total, filter)
	return resp, nil
}
