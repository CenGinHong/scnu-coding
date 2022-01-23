package define

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gmeta"
)

type ListCourseResp struct {
	gmeta.Meta    `orm:"table:course"`
	CourseId      int `orm:"course_id,primary" json:"courseId"` // 主键
	UserId        int `orm:"user_id" json:"-"`                  // 教师id
	TeacherDetail *struct {
		gmeta.Meta `orm:"table:sys_user"`
		UserId     int    `orm:"user_id" json:"userId"`    // 教师id
		Username   string `orm:"username" json:"username"` // 教师名
	} `orm:"with:user_id" json:"teacherDetail"`
	CourseName string      `orm:"course_name"       json:"courseName"` // 课程名称，限15字
	CourseDes  string      `orm:"course_des"        json:"courseDes"`  // 课程描述，限300字
	CoverImg   string      `orm:"cover_img"         json:"coverImg"`   // 封面url
	IsClose    bool        `orm:"is_close"          json:"isClose"`    // 结课标志
	CreatedAt  *gtime.Time `orm:"created_at"        json:"createdAt"`  // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"        json:"updatedAt"`  // 修改时间
}

type CourseDetailResp struct {
	gmeta.Meta    `orm:"table:course"`
	CourseId      int `orm:"course_id" json:"courseId"` // 主键
	UserId        int `orm:"user_id" json:"-"`          // 教师id
	TeacherDetail *struct {
		gmeta.Meta   `orm:"table:sys_user"`
		UserId       int    `orm:"user_id" json:"userId"`               // 教师id
		Username     string `orm:"username" json:"username"`            // 教师名
		CoverImg     string `orm:"cover_img"         json:"coverImg"`   // 封面url
		Organization string `orm:"organization"    json:"organization"` // 单位，例如计算机学院，限15字
		Email        string `orm:"email"           json:"email"`        // 邮箱，限30字
	} `orm:"with:user_id" json:"teacherDetail"`
	CourseName   string      `orm:"course_name" json:"courseName"`         // 课程名称，限15字
	CourseDes    string      `orm:"course_des" json:"courseDes"`           // 课程描述，限300字
	CoverImg     string      `orm:"cover_img"         json:"coverImg"`     // 封面url
	LanguageType int         `orm:"language_type"     json:"languageType"` //
	SecretKey    int         `orm:"secret_key" json:"secretKey"`           // 加入课程的密码,6位
	IsClose      bool        `orm:"is_close" json:"isClose"`               // 结课标志
	CreatedAt    *gtime.Time `orm:"created_at" json:"createdAt"`           // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at" json:"updatedAt"`           // 修改时间
}

type ConfirmStudentResp struct {
	UserId       int    `orm:"user_id,primary" json:"userId"`       // 主键
	Email        string `orm:"email"           json:"email"`        // 邮箱
	UserNum      string `orm:"user_num"        json:"userNum"`      // 学号/职工号
	Username     string `orm:"username"        json:"username"`     // 真实姓名
	Grade        int    `orm:"grade"           json:"grade"`        // 年级
	School       string `orm:"school"          json:"school"`       // 学院
	Gender       int    `orm:"gender"          json:"gender"`       // 性别
	Major        string `orm:"major"           json:"major"`        // 专业
	Organization string `orm:"organization"    json:"organization"` // 单位
}

type ImportStudent2Class struct {
	StudentNums []string
	CourseId    int
}

type ImportStudent2ClassResp struct {
	SuccessRecords   []*ConfirmStudentResp `json:"successRecords"`
	ErrorStudentNums []string              `json:"errorStudentNums"`
}

type SearchCourseResp struct {
	ListCourseResp
	IsTakeDetail struct {
		CourseId int  `orm:"course_id" json:"-"` // 主键
		IsTake   bool `json:"isTake"`            // 是否已经加入该门课程
	} `json:"isTakeDetail"`
}

type CourseStudentOverviewResp struct {
	gmeta.Meta `orm:"table:re_course_user"`
	UserId     int `orm:"user_id"         json:"userId"`
	UserDetail struct {
		gmeta.Meta   `orm:"table:sys_user"`
		UserId       int    `orm:"user_id"         json:"-"`            // 用户id
		Email        string `orm:"email"           json:"email"`        // 邮箱，限30字
		UserNum      string `orm:"user_num"        json:"userNum"`      // 学号/职工号，限20位
		Grade        int    `orm:"grade"           json:"grade"`        // 年级
		School       string `orm:"school"          json:"school"`       // 学院
		Major        string `orm:"major"           json:"major"`        // 专业
		Username     string `orm:"username"        json:"username"`     // 真实姓名，限6字
		Organization string `orm:"organization"    json:"organization"` // 单位，例如计算机学院，限15字
	} `orm:"with:user_id" json:"userDetail"`
	AvgScoreDetail struct {
		Score float32 `json:"score"`
	} `json:"avgScoreDetail"`
	CheckinDetail struct {
		CheckinCount int `json:"checkinCount"`
		TotalCount   int `json:"totalCount"`
	} `json:"checkinDetail"`
	CodingTimeDetail struct {
		CodingTime int `json:"codingTime"`
	} `json:"codingTimeDetail"`
}

type DeleteStudentFromClassReq struct {
	UserIds  []int
	CourseId int
}

type InsertCourseReq struct {
	CourseName   string `orm:"course_name"       json:"courseName"`   // 课程名称，限15字
	UserId       int    `orm:"user_id"           json:"userId"`       // 教师id
	CourseDes    string `orm:"course_des"        json:"courseDes"`    // 课程描述，限300字
	CoverImg     string `orm:"cover_img"         json:"coverImg"`     // 封面url
	SecretKey    int    `orm:"secret_key"        json:"secretKey"`    // 加入课程的密码,6位
	LanguageType int    `orm:"language_type"     json:"languageType"` // 语言类型
}

type UpdateCourseReq struct {
	CourseId   int     `orm:"course_id"`   // 主键
	CourseName string  `orm:"course_name"` // 课程名称，限15字
	CourseDes  string  `orm:"course_des"`  // 课程描述，限300字
	CoverImg   *string `orm:"cover_img"`   // 封面url
	SecretKey  int     `orm:"secret_key"`  // 加入课程的密码,6位
	IsClose    bool    `orm:"is_close"`
}

type EnrollUserDetail struct {
	gmeta.Meta `orm:"table:re_course_user"`
	UserId     int `orm:"user_id"         json:"-"`
	UserDetail struct {
		gmeta.Meta   `orm:"table:sys_user"`
		UserId       int    `orm:"user_id"         json:"userId"`       // 用户id
		Email        string `orm:"email"           json:"email"`        // 邮箱，限30字
		UserNum      string `orm:"user_num"        json:"userNum"`      // 学号/职工号，限20位
		Grade        int    `orm:"grade"           json:"grade"`        // 年级
		School       string `orm:"school"          json:"school"`       // 学院
		Major        string `orm:"major"           json:"major"`        // 专业
		Username     string `orm:"username"        json:"username"`     // 真实姓名，限6字
		Organization string `orm:"organization"    json:"organization"` // 单位，例如计算机学院，限15字
	} `orm:"with:user_id" json:"userDetail"`
}

type JoinClassReq struct {
	UserId    int
	CourseId  int
	SecretKey string
}
