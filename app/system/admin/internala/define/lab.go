package define

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gmeta"
)

type Lab struct {
	gmeta.Meta   `orm:"table:lab"`
	LabId        int `orm:"lab_id,primary" json:"labId"`    // 主键
	CourseId     int `orm:"course_id"      json:"courseId"` // 该实验隶属的课程
	CourseDetail *struct {
		gmeta.Meta `orm:"table:course"`
		CourseId   int    `orm:"course_id,primary" json:"-"`          // 主键
		UserId     int    `orm:"user_id"           json:"userId"`     // 教师id
		CourseName string `orm:"course_name"       json:"courseName"` // 课程名称，限15字
	} `orm:"with:course_id" json:"courseDetail"`
	Type          int         `orm:"type"           json:"type"`          // 枚举，练习，作业，考试
	Title         string      `orm:"title"          json:"title"`         // 标题
	Content       string      `orm:"content"        json:"content"`       // 实验内容描述
	CreatedAt     *gtime.Time `orm:"created_at"     json:"createdAt"`     // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"     json:"updatedAt"`     // 修改时间
	AttachmentSrc string      `orm:"attachment_src" json:"attachmentSrc"` // 实验附件url
	Deadline      string      `orm:"deadline"       json:"deadline"`      // 截止时间
}
