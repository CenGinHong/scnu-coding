package define

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gmeta"
)

type Course struct {
	gmeta.Meta    `orm:"table:course"`
	CourseId      int `orm:"course_id,primary" json:"courseId"` // 主键
	UserId        int `orm:"user_id"           json:"-"`        // 教师id
	TeacherDetail struct {
		gmeta.Meta `orm:"table:sys_user"`
		UserId     int    `orm:"user_id"           json:"userId"`    // 教师id
		UserName   string `orm:"userName"           json:"username"` // 教师名称
	} `orm:"with:user_id" json:"teacherDetail"`
	EnrollCount  int         `json:"enrollCount"`                          // 选课总人数
	CourseName   string      `orm:"course_name"       json:"courseName"`   // 课程名称，限15字
	CourseDes    string      `orm:"course_des"        json:"courseDes"`    // 课程描述，限300字
	CoverUrl     string      `orm:"cover_url"         json:"coverUrl"`     // 封面url
	SecretKey    int         `orm:"secret_key"        json:"secretKey"`    // 加入课程的密码,6位
	IsClose      int         `orm:"is_close"          json:"isClose"`      // 结课标志
	LanguageType int         `orm:"language_type"     json:"languageType"` //
	CreatedAt    *gtime.Time `orm:"created_at"        json:"createdAt"`    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"        json:"updatedAt"`    // 修改时间
}
