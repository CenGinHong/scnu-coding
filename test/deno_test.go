package test

import (
	"github.com/gogf/gf/os/gtime"
	"scnu-coding/app/dao"
	"scnu-coding/app/model"
	"testing"
)

type demo struct {
	UserId int
	Detail *struct {
		UserId   int
		UserName string
	}
}

type ListCourseByTeacherIdResp struct {
	model.Course
}

type CodingTimeRecord struct {
	Duration  int         `orm:"duration"               json:"duration"`                                        // 编码时间，分钟为单位
	CreatedAt *gtime.Time `orm:"Date_Format(created_at,'%Y-%m-%d') as created_at"             json:"createdAt"` // 创建时间
}

func TestGetAllUser(t *testing.T) {
	demo := make([]*ListCourseByTeacherIdResp, 0)
	err := dao.Course.Scan(&demo)
	if err != nil {
		return
	}
	println(demo)
}
