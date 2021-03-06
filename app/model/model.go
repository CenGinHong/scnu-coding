// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package model

import (
	"github.com/gogf/gf/os/gtime"
)

// CheckinDetail is the golang structure for table checkin_detail.
type CheckinDetail struct {
	CheckinDetailId int         `orm:"checkin_detail_id,primary" json:"checkinDetailId"` // id
	IsCheckin       int         `orm:"is_checkin"                json:"isCheckin"`       // 是否有签到
	UserId          int         `orm:"user_id"                   json:"userId"`          // 用户id
	CheckinRecordId int         `orm:"checkin_record_id"         json:"checkinRecordId"` // 签到记录id
	CreatedAt       *gtime.Time `orm:"created_at"                json:"createdAt"`       // 创建时间
	UpdatedAt       *gtime.Time `orm:"updated_at"                json:"updatedAt"`       // 更新时间
}

// CheckinRecord is the golang structure for table checkin_record.
type CheckinRecord struct {
	CheckinRecordId int         `orm:"checkin_record_id,primary" json:"checkinRecordId"` // id
	CheckinName     string      `orm:"checkin_name"              json:"checkinName"`     // 签到名称，例如2021年2月5日签到
	CheckinKey      string      `orm:"checkin_key"               json:"checkinKey"`      // 签到密钥
	TotalDuration   int         `orm:"total_duration"            json:"totalDuration"`   // 限时时间
	CourseId        int         `orm:"course_id"                 json:"courseId"`        // 课程id
	CreatedAt       *gtime.Time `orm:"created_at"                json:"createdAt"`       // 创建时间
	UpdatedAt       *gtime.Time `orm:"updated_at"                json:"updatedAt"`       // 更新时间
}

// CodingTime is the golang structure for table coding_time.
type CodingTime struct {
	CodingTimeId int         `orm:"coding_time_id,primary" json:"codingTimeId"` //
	LabId        int         `orm:"lab_id"                 json:"labId"`        // 实验id
	UserId       int         `orm:"user_id"                json:"userId"`       // 学生id
	Duration     int         `orm:"duration"               json:"duration"`     // 编码时间，分钟为单位
	CreatedAt    *gtime.Time `orm:"created_at"             json:"createdAt"`    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"             json:"updatedAt"`    // 更新时间
}

// Course is the golang structure for table course.
type Course struct {
	CourseId     int         `orm:"course_id,primary" json:"courseId"`     // 主键
	UserId       int         `orm:"user_id"           json:"userId"`       // 教师id
	CourseName   string      `orm:"course_name"       json:"courseName"`   // 课程名称，限15字
	CourseDes    string      `orm:"course_des"        json:"courseDes"`    // 课程描述，限300字
	CoverImg     string      `orm:"cover_img"         json:"coverImg"`     // 封面url
	SecretKey    int         `orm:"secret_key"        json:"secretKey"`    // 加入课程的密码,6位
	IsClose      int         `orm:"is_close"          json:"isClose"`      // 结课标志
	LanguageType int         `orm:"language_type"     json:"languageType"` //
	CreatedAt    *gtime.Time `orm:"created_at"        json:"createdAt"`    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"        json:"updatedAt"`    // 修改时间
}

// CourseAnnouncement is the golang structure for table course_announcement.
type CourseAnnouncement struct {
	CourseAnnouncementId int         `orm:"course_announcement_id,primary" json:"courseAnnouncementId"` // id
	Title                string      `orm:"title"                          json:"title"`                // 标题
	CourseId             int         `orm:"course_id"                      json:"courseId"`             // 课程id
	Content              string      `orm:"content"                        json:"content"`              // 公告内容，限2000字
	AttachmentSrc        string      `orm:"attachment_src"                 json:"attachmentSrc"`        // 文件url
	CreatedAt            *gtime.Time `orm:"created_at"                     json:"createdAt"`            // 创建时间
	UpdatedAt            *gtime.Time `orm:"updated_at"                     json:"updatedAt"`            // 修改时间
}

// CourseComment is the golang structure for table course_comment.
type CourseComment struct {
	CourseCommentId int         `orm:"course_comment_id,primary" json:"courseCommentId"` // 主键
	CourseId        int         `orm:"course_id"                 json:"courseId"`        // 实验内容
	CommentText     string      `orm:"comment_text"              json:"commentText"`     // 评论内容，限120字
	Pid             int         `orm:"pid"                       json:"pid"`             // 父评论id，主评时为空
	UserId          int         `orm:"user_id"                   json:"userId"`          // 发评论的用户id
	CreatedAt       *gtime.Time `orm:"created_at"                json:"createdAt"`       // 创建时间
	UpdatedAt       *gtime.Time `orm:"updated_at"                json:"updatedAt"`       // 更新时间
}

// Lab is the golang structure for table lab.
type Lab struct {
	LabId         int         `orm:"lab_id,primary" json:"labId"`         // 主键
	CourseId      int         `orm:"course_id"      json:"courseId"`      // 该实验隶属的课程
	Type          int         `orm:"type"           json:"type"`          // 枚举，练习，作业，考试
	Title         string      `orm:"title"          json:"title"`         // 标题
	Content       string      `orm:"content"        json:"content"`       // 实验内容描述
	CreatedAt     *gtime.Time `orm:"created_at"     json:"createdAt"`     // 创建时间
	UpdatedAt     *gtime.Time `orm:"updated_at"     json:"updatedAt"`     // 修改时间
	AttachmentSrc string      `orm:"attachment_src" json:"attachmentSrc"` // 实验附件url
	Deadline      string      `orm:"deadline"       json:"deadline"`      // 截止时间
}

// LabComment is the golang structure for table lab_comment.
type LabComment struct {
	LabCommentId int         `orm:"lab_comment_id,primary" json:"labCommentId"` // 主键
	LabId        int         `orm:"lab_id"                 json:"labId"`        // 实验id
	CommentText  string      `orm:"comment_text"           json:"commentText"`  // 评论内容，限120字
	Pid          int         `orm:"pid"                    json:"pid"`          // 父评论id，主评时为0
	UserId       int         `orm:"user_id"                json:"userId"`       // 发评论的用户id
	CreatedAt    *gtime.Time `orm:"created_at"             json:"createdAt"`    // 创建时间
	UpdatedAt    *gtime.Time `orm:"updated_at"             json:"updatedAt"`    // 更新时间
}

// LabSubmit is the golang structure for table lab_submit.
type LabSubmit struct {
	LabSubmitId      int         `orm:"lab_submit_id,primary" json:"labSubmitId"`      //
	LabId            int         `orm:"lab_id"                json:"labId"`            // lab id
	UserId           int         `orm:"user_id"               json:"userId"`           // 用户id
	ReportContent    string      `orm:"report_content"        json:"reportContent"`    // 实验报告(md)/存放实验报告pdf的url
	Score            int         `orm:"score"                 json:"score"`            // 分数
	IsFinish         int         `orm:"is_finish"             json:"isFinish"`         // 是否完成
	LabSubmitComment string      `orm:"lab_submit_comment"    json:"labSubmitComment"` // 评论
	UpdatedAt        *gtime.Time `orm:"updated_at"            json:"updatedAt"`        // 更新时间
	CreatedAt        *gtime.Time `orm:"created_at"            json:"createdAt"`        // 创建时间
}

// Log is the golang structure for table log.
type Log struct {
	LogId      int    `orm:"log_id,primary" json:"logId"`      //
	UserId     int    `orm:"user_id"        json:"userId"`     //
	HttpMethod string `orm:"http_method"    json:"httpMethod"` // 操作类型（GET,POST,PUT,DELETE)
	RespData   string `orm:"resp_data"      json:"respData"`   //
	ReqParams  string `orm:"req_params"     json:"reqParams"`  //
}

// MessageNotify is the golang structure for table message_notify.
type MessageNotify struct {
	MessageId  int         `orm:"message_id,primary" json:"messageId"`  // id
	ReceiverId int         `orm:"receiver_id"        json:"receiverId"` // 接收者id
	SenderId   int         `orm:"sender_id"          json:"senderId"`   // 发送者id
	Content    string      `orm:"content"            json:"content"`    // 消息内容
	CreatedAt  *gtime.Time `orm:"created_at"         json:"createdAt"`  // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"         json:"updatedAt"`  // 更新时间
}

// ReApiRole is the golang structure for table re_api_role.
type ReApiRole struct {
	ApiRoleId int         `orm:"api_role_id,primary" json:"apiRoleId"` //
	RoleId    int         `orm:"role_id"             json:"roleId"`    //
	ApiId     int         `orm:"api_id"              json:"apiId"`     //
	CreatedAt *gtime.Time `orm:"created_at"          json:"createdAt"` // 创建时间
	UpdatedAt *gtime.Time `orm:"updated_at"          json:"updatedAt"` // 修改时间
	DeletedAt *gtime.Time `orm:"deleted_at"          json:"deletedAt"` // 删除标记
}

// ReCourseUser is the golang structure for table re_course_user.
type ReCourseUser struct {
	ReCourseUserId int         `orm:"re_course_user_id,primary" json:"reCourseUserId"` // 主键
	UserId         int         `orm:"user_id"                   json:"userId"`         // 用户id
	CourseId       int         `orm:"course_id"                 json:"courseId"`       // 课程id
	UpdatedAt      *gtime.Time `orm:"updated_at"                json:"updatedAt"`      // 更新时间
	CreatedAt      *gtime.Time `orm:"created_at"                json:"createdAt"`      // 创建时间
}

// ReUserRole is the golang structure for table re_user_role.
type ReUserRole struct {
	UserRoleId int         `orm:"user_role_id,primary" json:"userRoleId"` // 主键
	UserId     int         `orm:"user_id"              json:"userId"`     // 用户id
	RoleId     int         `orm:"role_id"              json:"roleId"`     // 角色id
	CreatedAt  *gtime.Time `orm:"created_at"           json:"createdAt"`  // 创建时间
	UpdatedAt  *gtime.Time `orm:"updated_at"           json:"updatedAt"`  // 修改时间
}

// SysApi is the golang structure for table sys_api.
type SysApi struct {
	ApiId       int         `orm:"api_id,primary" json:"apiId"`       //
	Api         string      `orm:"api"            json:"api"`         // api
	Method      string      `orm:"method"         json:"method"`      //
	Description string      `orm:"description"    json:"description"` // 描述
	CreatedAt   *gtime.Time `orm:"created_at"     json:"createdAt"`   // 创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at"     json:"updatedAt"`   // 修改时间
}

// SysRole is the golang structure for table sys_role.
type SysRole struct {
	RoleId      int         `orm:"role_id,primary" json:"roleId"`      // 主键
	Description string      `orm:"description"     json:"description"` // 描述
	CreatedAt   *gtime.Time `orm:"created_at"      json:"createdAt"`   // 创建时间
	UpdatedAt   *gtime.Time `orm:"updated_at"      json:"updatedAt"`   // 更新时间
}

// SysUser is the golang structure for table sys_user.
type SysUser struct {
	UserId       int         `orm:"user_id,primary" json:"userId"`       // 主键
	Email        string      `orm:"email"           json:"email"`        // 邮箱
	UserNum      string      `orm:"user_num"        json:"userNum"`      // 学号/职工号
	Username     string      `orm:"username"        json:"username"`     // 真实姓名
	Grade        int         `orm:"grade"           json:"grade"`        // 年级
	School       string      `orm:"school"          json:"school"`       // 学院
	Password     string      `orm:"password"        json:"password"`     // 密码
	AvatarImg    string      `orm:"avatar_img"      json:"avatarImg"`    // 头像url
	Gender       int         `orm:"gender"          json:"gender"`       // 性别
	Major        string      `orm:"major"           json:"major"`        // 专业
	Organization string      `orm:"organization"    json:"organization"` // 单位，例如软件工程4班
	RoleId       int         `orm:"role_id"         json:"roleId"`       // 角色id
	UpdatedAt    *gtime.Time `orm:"updated_at"      json:"updatedAt"`    // 修改时间
	CreatedAt    *gtime.Time `orm:"created_at"      json:"createdAt"`    // 创建时间
}
