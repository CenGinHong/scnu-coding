package define

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gmeta"
)

type ListCourseResp struct {
	gmeta.Meta    `orm:"table:course"`
	CourseId      int `orm:"course_id,primary" json:"courseId"` // 主键
	UserId        int `orm:"user_id"           json:"userId"`   // 教师id
	TeacherDetail *struct {
		gmeta.Meta `orm:"table:sys_user"`
		UserId     int    `orm:"user_id" json:"-"`         // 教师id
		Username   string `orm:"username" json:"username"` // 教师名
	} `orm:"with:user_id" json:"teacherDetail"`
	CourseName   string      `orm:"course_name"       json:"courseName"`   // 课程名称，限15字
	CourseDes    string      `orm:"course_des"        json:"courseDes"`    // 课程描述，限300字
	CoverImg     string      `orm:"cover_img"         json:"coverImg"`     // 封面url
	SecretKey    int         `orm:"secret_key"        json:"secretKey"`    // 加入课程的密码,6位
	IsClose      int         `orm:"is_close"          json:"isClose"`      // 结课标志
	LanguageType int         `orm:"language_type"     json:"languageType"` //
	UpdatedAt    *gtime.Time `orm:"updated_at"        json:"updatedAt"`    // 修改时间
}

// CourseEnroll 选课的学生
type CourseEnroll struct {
	UserId       int    `orm:"a.user_id"       json:"userId"`       // 用户id
	Email        string `orm:"email"           json:"email"`        // 邮箱，限30字
	UserNum      string `orm:"user_num"        json:"userNum"`      // 学号/职工号，限20位
	Major        string `orm:"major"           json:"major"`        // 专业
	Username     string `orm:"username"        json:"username"`     // 真实姓名，限6字
	Organization string `orm:"organization"    json:"organization"` // 单位，例如计算机学院，限15字
}

type RemoveCourseEnrollReq struct {
	UserIds  []int
	CourseId int
}

type AddStudent2ClassReq struct {
	CourseId    int
	StudentNums []string
}
