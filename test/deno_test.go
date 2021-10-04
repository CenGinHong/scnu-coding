package test

import (
	"context"
	"fmt"
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

type CourseEnroll struct {
	UserId       int    `orm:"user_id"         json:"userId"`
	Email        string `orm:"email"           json:"email"`        // 邮箱，限30字
	UserNum      string `orm:"user_num"        json:"userNum"`      // 学号/职工号，限20位
	Grade        int    `orm:"grade"           json:"grade"`        // 年级
	School       string `orm:"school"          json:"school"`       // 学院
	Gender       int    `orm:"gender"          json:"gender"`       // 性别
	Major        string `orm:"major"           json:"major"`        // 专业
	Username     string `orm:"username"        json:"username"`     // 真实姓名，限6字
	Organization string `orm:"organization"    json:"organization"` // 单位，例如计算机学院，限15字
}

type CodingTimeRecord struct {
	Duration  int         `orm:"duration"               json:"duration"`                                        // 编码时间，分钟为单位
	CreatedAt *gtime.Time `orm:"Date_Format(created_at,'%Y-%m-%d') as created_at"             json:"createdAt"` // 创建时间
}

func TestGetAllUser(t *testing.T) {
	// 获取分页信息
	courseId := 126
	records := make([]*CourseEnroll, 0)
	tableAlias := "a"
	d := dao.ReCourseUser.Ctx(context.Background()).As(tableAlias)
	d = d.LeftJoin(dao.SysUser.Table, fmt.Sprintf("%s.%s=%s.%s", tableAlias, dao.ReCourseUser.Columns.UserId, dao.SysUser.Table, dao.SysUser.Columns.UserId))
	d = d.Where(fmt.Sprintf("%s.%s", tableAlias, dao.ReCourseUser.Columns.CourseId), courseId)
	err := d.Scan(&records)
	if err != nil {
		return
	}

}
