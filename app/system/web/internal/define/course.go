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
	CourseName      string `orm:"course_name"       json:"courseName"` // 课程名称，限15字
	CourseDes       string `orm:"course_des"        json:"courseDes"`  // 课程描述，限300字
	CoverFileId     string `orm:"cover_file_id"     json:"-"`          // 封面url的id
	CoverFileDetail *struct {
		gmeta.Meta  `orm:"table:local_file"`
		LocalFileId int    `orm:"local_file_id,primary" json:"-"`   //
		Url         string `orm:"url"                   json:"url"` // 文件url
	} `orm:"with:local_file_id=cover_file_id" json:"coverFileDetail"`
	IsClose   bool        `orm:"is_close"          json:"isClose"`   // 结课标志
	CreatedAt *gtime.Time `orm:"created_at"        json:"createdAt"` // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at"        json:"updatedAt"` // 修改时间
}

type CourseDetailResp struct {
	gmeta.Meta    `orm:"table:course"`
	CourseId      int `orm:"course_id" json:"courseId"` // 主键
	UserId        int `orm:"user_id" json:"-"`          // 教师id
	TeacherDetail *struct {
		gmeta.Meta       `orm:"table:sys_user"`
		UserId           int    `orm:"user_id" json:"userId"`    // 教师id
		Username         string `orm:"username" json:"username"` // 教师名
		AvatarFileId     string `orm:"avatar_file_id"  json:"-"` // 头像url
		AvatarFileDetail *struct {
			gmeta.Meta  `orm:"table:local_file"`
			LocalFileId int    `orm:"local_file_id,primary" json:"-"`   //
			Url         string `orm:"url"                   json:"url"` // 文件url
		} `orm:"with:local_file_id=avatar_file_id" json:"avatarFileDetail"`
		Organization string `orm:"organization"    json:"organization"` // 单位，例如计算机学院，限15字
		Email        string `orm:"email"           json:"email"`        // 邮箱，限30字
	} `orm:"with:user_id" json:"teacherDetail"`
	CourseName      string `orm:"course_name" json:"courseName"` // 课程名称，限15字
	CourseDes       string `orm:"course_des" json:"courseDes"`   // 课程描述，限300字
	CoverFileId     string `orm:"cover_file_id"     json:"-"`    // 封面url的id
	CoverFileDetail *struct {
		gmeta.Meta  `orm:"table:local_file"`
		LocalFileId int    `orm:"local_file_id,primary" json:"-"`   //
		Url         string `orm:"url"                   json:"url"` // 文件url
	} `orm:"with:local_file_id=cover_file_id" json:"coverFileDetail"`
	SecretKey int         `orm:"secret_key" json:"secretKey"` // 加入课程的密码,6位
	IsClose   bool        `orm:"is_close" json:"isClose"`     // 结课标志
	CreatedAt *gtime.Time `orm:"created_at" json:"createdAt"` // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at" json:"updatedAt"` // 修改时间
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

//
//type EnrollCourseReq struct {
//	StudentId int // 学生id
//	CourseId  int // 课程id
//	SecretKey int // 密钥
//}
//
//type UpdateCourseReq struct {
//	CourseId   int     `orm:"course_id"`   // 主键
//	CourseName *string `orm:"course_name"` // 课程名称，限长15字
//	CourseDes  *string `orm:"course_des"`  // 课程描述,限两百字
//	CoverURL   *string `orm:"cover_url"`   // 封面url
//	SecretKey  *int    `orm:"secret_key"`  // 加入课程的密码,6位
//}
//
//type InsertCourseReq struct {
//	CourseName string `orm:"course_name"` // 课程名称，限长15字
//	CourseDes  string `orm:"course_des"`  // 课程描述,限两百字
//	CoverURL   string `orm:"cover_url"`   // 封面url,可空
//	SecretKey  int    `orm:"secret_key"`  // 加入课程的密码,6位
//	Language   int    `orm:"language"`    // 语言类型枚举
//}
//
//type ListCodingTimeByCourseIdResp struct {
//	UserId     int    `orm:"user_id" json:"userId"`
//	Username   string `orm:"username" json:"username"`
//	UserNum    string `orm:"user_num" json:"userNum"`
//	CodingTime []*CodingTimeRecord
//}

//type CourseScore struct {
//	gmeta.Meta `orm:"table:re_course_user"`
//	StuID     int `orm:"user_id" json:"userId"`
//	CourseId   int `orm:"course_id" json:"courseId"`
//	UserDetail *struct {
//		gmeta.Meta `orm:"table:sys_user"`
//		StuID     int    `orm:"user_id" json:"userId"`    // 学生id
//		Username   string `orm:"username" json:"username"` // 学生名
//		UserNum    string `orm:"user_num" json:"userNum"`  // 学号
//	} `orm:"with:user_id"` // 学生详细信息
//	ScoreDetail []*struct {
//		gmeta.Meta `orm:"table:lab"`
//		LabID      int    `orm:"lab_id"`    // labId
//		Title      string `orm:"title"`     // 实验名
//		CourseId   int    `orm:"course_id"` // 课程id
//		LabDetail  *struct {
//		} `orm:"with:course_id" json:"labDetail"` // 该课程内的所有实验成绩
//	} `orm:"with:course_id"`               // 该课程内的所有实验
//	AvgScore     int `orm:"avg_score"`     // 平均成绩
//	ShallCheckIn int `orm:"shall_checkin"` // 应该签到的次数
//	ActCheckIn   int `orm:"act_checkin"`   // 实际签到的次数
//}

//type ListCourseScoreResp struct {
//	StuID       int    // 学生id
//	Num          string // 学生学号
//	RealName     string // 学生姓名
//	AvgScore     int    // 平均成绩
//	ShallCheckIn int    // 应该签到的次数
//	ActCheckIn   int    // 实际签到的次数
//}

//type ExamineEnrollReq struct {
//	StuIds      []int // 允许加入课程
//	CourseId    int   // 课程学生
//	IsPermitted bool  // 是否允许加入课程
//}
//
//type DropCourseReq struct {
//	CourseId int   // 课程id
//	StuIds   []int // 学生id
//}
//
//type ListStuWaitExaminedResp struct {
//	gmeta.Meta `orm:"table:re_course_user"`
//	UserId     int `orm:"user_id" json:"userId"` // 用户id
//	UserDetail *struct {
//		gmeta.Meta `orm:"table:sys_user"`
//		UserId     int    `orm:"user_id" json:"userId"`    // 学生id
//		Username   string `orm:"username" json:"username"` // 学生名
//		UserNum    string `orm:"user_num" json:"userNum"`  // 学号
//	} `orm:"with:user_id" json:"userDetail"` // 学生详细信息
//}
//
//type CourseScoreCsv struct {
//	gmeta.Meta `orm:"table:lab_submit"`
//	LabId      int `orm:"lab_id"  json:"labId"`  // lab id
//	UserId     int `orm:"user_id" json:"userId"` // 用户id
//	UserDetail *struct {
//		gmeta.Meta `orm:"table:sys_user"`
//		UserId     int    `orm:"user_id" json:"userId"`    // 用户id
//		Username   string `orm:"username" json:"username"` //
//		UserNum    string `orm:"user_num" json:"userNum"`  // 学号
//	} `orm:"with:user_id" json:"userDetail"`
//	Score *int `orm:"score" json:"score"` // 分数,这里要用指针，区分未评分和已评分
//}
//
////================================================================================

type SearchCourseResp struct {
	ListCourseResp
	IsTakeDetail struct {
		IsTake bool `json:"isTake"` // 是否已经加入该门课程
	} `json:"isTakeDetail"`
}

type CourseStudentOverviewResp struct {
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
