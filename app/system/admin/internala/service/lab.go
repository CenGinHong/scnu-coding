package service

import (
	"context"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/model"
	"scnu-coding/app/service"
	"scnu-coding/app/system/admin/internala/define"
	"scnu-coding/library/response"
)

var Lab = labService{}

type labService struct {
}

func (s *labService) ListAllLab(ctx context.Context) (resp *response.PageResp, err error) {
	// 获取分页信息
	pageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.Lab, 0)
	d := dao.Lab.Ctx(ctx)
	if pageInfo != nil {
		// 筛选条件
		for key, value := range pageInfo.ParseFilterFields {
			d = d.Where(key, value)
		}
		// 升降序
		if pageInfo.SortOrder != "" {
			d = d.Order(pageInfo.SortField, pageInfo.SortOrder)
		}
		d = d.Page(pageInfo.Current, pageInfo.PageSize)
	}
	total, err := d.Count()
	if err != nil {
		return nil, nil
	}
	if err = d.WithAll().Scan(&records); err != nil {
		return nil, err
	}
	// 构建筛选集
	courseModels := make([]model.Course, 0)
	if err = dao.Course.Ctx(ctx).Fields(dao.Course.Columns.CourseId, dao.Course.Columns.CourseName).Scan(&courseModels); err != nil {
		return nil, err
	}
	filter := make(map[string][]*response.FilterType, 0)
	tempFilter := make([]*response.FilterType, 0)
	for _, courseFilter := range courseModels {
		tempFilter = append(tempFilter, &response.FilterType{
			Text:     courseFilter.CourseName,
			Value:    gconv.String(courseFilter.CourseId),
			Children: nil,
		})
	}
	filter["courseName"] = tempFilter
	resp = response.GetPageResp(records, total, filter)
	return resp, nil
}
