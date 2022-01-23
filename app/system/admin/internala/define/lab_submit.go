package define

import (
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gmeta"
)

type ListLabSubmitByLabIdResp struct {
	gmeta.Meta `orm:"table:re_course_user"`
	UserId     int `orm:"user_id" json:"userId"` // 用户id
	UserDetail *struct {
		gmeta.Meta `orm:"table:sys_user"`
		UserId     int    `orm:"user_id" json:"-"`         // 用户id
		Username   string `orm:"username" json:"username"` //
		UserNum    string `orm:"user_num" json:"userNum"`  // 学号
	} `orm:"with:user_id" json:"userDetail"`
	LabSubmitDetail struct {
		LabSubmitId      int         `orm:"lab_submit_id" json:"labSubmitId"` //
		IsFinish         bool        `orm:"is_finish" json:"isFinish"`
		Score            int         `orm:"score" json:"score"`
		UserId           int         `orm:"user_id" json:"-"`
		LabSubmitComment string      `orm:"lab_submit_comment" json:"labSubmitComment"` // 评语
		UpdatedAt        *gtime.Time `orm:"updated_at" json:"updatedAt"`                // 更新时间
	} `json:"labSubmitDetail"`
}
