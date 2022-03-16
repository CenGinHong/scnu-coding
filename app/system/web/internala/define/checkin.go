package define

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gmeta"
)

// CheckinStatusResp 获取正在课程签到的实时信息,该结构是RedisCheckinData的子集
type CheckinStatusResp struct {
	CheckinRecordId int     `orm:"checkin_record_id" json:"checkinRecordId"` // 签到id
	CheckinName     string  `orm:"checkin_name" json:"checkinName"`
	RemainDuration  float64 `json:"remainDuration"`
	TotalDuration   float64 `orm:"total_duration" json:"totalDuration"`
}

type RedisCheckinData struct {
	CheckinRecordId int64   // id
	CheckinName     string  // 签到名称
	CheckinKey      string  // 密钥
	TotalDuration   float64 // 总签到允许时间
}

type StartCheckInReq struct {
	CheckinName string
	CheckinKey  string
	CourseId    int
	Duration    float64
}

type UpdateCheckinDetailReq struct {
	UserId          int
	CheckinRecordId int
	IsCheckin       bool
}

type StudentCheckinReq struct {
	CourseId   int    // 签到id
	CheckinKey string // 课程key
}

type CheckinDetailResp struct {
	gmeta.Meta `orm:"table:re_course_user"`
	UserId     int `orm:"user_id" json:"userId"` // 学生id
	UserDetail *struct {
		gmeta.Meta `orm:"table:sys_user"`
		UserId     int    `orm:"user_id" json:"-"`         // 主键
		UserNum    string `orm:"user_num" json:"userNum"`  // 学号/职工号，限20位
		Username   string `orm:"username" json:"username"` // 真实姓名，限10字
	} `orm:"with:user_id" json:"userDetail"`
	CheckinDetail struct {
		UserId    int  `orm:"user_id" json:"-"`            // 主键
		IsCheckin bool `orm:"is_checkin" json:"isCheckin"` // 是否已经签到
	} `json:"checkinDetail"`
}

type StuListCheckInRecordResp struct {
	CheckinRecordId int    `orm:"checkin_record_id,primary" json:"checkinRecordId"` // id
	CheckinName     string `orm:"checkin_name" json:"checkinName"`                  // 签到名称，例如2021年2月5日签到
	CheckinDetail   struct {
		CheckinRecordId int  `orm:"checkin_record_id" json:"-"`  // 签到记录id
		IsCheckin       bool `orm:"is_checkin" json:"isCheckin"` // 是否签到
	} `json:"checkinDetail"`
	CreatedAt *gtime.Time `orm:"created_at" json:"createdAt"` // 创建时间
}

type CheckinRecordResp struct {
	CheckinRecordId int         `orm:"checkin_record_id,primary" json:"checkinRecordId"` // id
	CheckinName     string      `orm:"checkin_name"              json:"checkinName"`     // 签到名称，例如2021年2月5日签到
	CheckinKey      string      `orm:"checkin_key"               json:"checkinKey"`      // 签到密钥
	TotalDuration   int         `orm:"total_duration"            json:"totalDuration"`   // 限时时间
	CreatedAt       *gtime.Time `orm:"created_at"                json:"createdAt"`       // 创建时间
	Attendance      struct {
		CheckinCount int `json:"checkinCount"`
		TakeCount    int `json:"takeCount"`
	} `json:"attendance"`
}

type ExportCheckinRecord struct {
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
	CheckinDetails []*struct {
		CheckinRecordId int  // 签到记录id
		UserId          int  // 参与签到的人
		IsCheckin       bool // 是否签到
	}
}
