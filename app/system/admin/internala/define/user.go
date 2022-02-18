package define

import (
	"github.com/gogf/gf/util/gmeta"
	"scnu-coding/app/model"
)

type SysUserResp struct {
	gmeta.Meta `orm:"table:sys_user"`
	model.SysUser
	//UserId       int    `orm:"user_id,primary" json:"userId"`       // 主键
	//Email        string `orm:"email"           json:"email"`        // 邮箱，限30字
	//Grade        int    `orm:"grade"           json:"grade"`        // 年级
	//UserNum      string `orm:"user_num"        json:"userNum"`      // 学号/职工号，限20位
	//Username     string `orm:"username"        json:"username"`     // 真实姓名，限6字
	//Gender       int    `orm:"gender"          json:"gender"`       // 性别
	//Major        string `orm:"major"           json:"major"`        // 专业，限15字
	//School       string `orm:"school"          json:"school"`       // 学院
	//Organization string `orm:"organization"    json:"organization"` // 单位，例如计算机学院，限15字
	//RoleId       int    `orm:"role_id"         json:"-"`            // 角色id
	RoleDetail struct {
		gmeta.Meta  `orm:"table:sys_role"`
		RoleId      int    `orm:"role_id,primary" json:"-"`
		Description string `orm:"description" json:"description"`
	} `orm:"with:role_id" json:"roleDetail"` // 角色信息
	//UpdatedAt *gtime.Time `orm:"updated_at"      json:"updatedAt"` // 修改时间
	//CreatedAt *gtime.Time `orm:"created_at"      json:"createdAt"` // 创建时间
}

type UpdateSysUserReq struct {
	UserId       int    `orm:"user_id,primary"` // 主键
	Email        string `orm:"email"`           // 邮箱，限30字
	UserNum      string `orm:"user_num"`        // 学号/职工号，限20位
	Username     string `orm:"username"`        // 真实姓名，限6字
	Gender       int    `orm:"gender"`          // 性别
	Major        string `orm:"major"`           // 专业，限15字
	Organization string `orm:"organization"`    // 单位，例如计算机学院，限15字
}

type ImportStudent struct {
	Username     string `orm:"username"        json:"username"     valid:"required#姓名不能为空" csv:"姓名"`                  // 真实姓名，限6字
	UserNum      string `orm:"user_num"        json:"userNum"      valid:"required|integer#学号不能为空|学号格式不正确"  csv:"学号"` // 学号/职工号，限20位
	School       string `orm:"school"          json:"school"  csv:"学院"`                                               // 学院
	Email        string `orm:"email"           json:"email"        valid:"email#邮箱字段不符合格式"  csv:"邮箱"`                 // 邮箱，限30字
	Organization string `orm:"organization"    json:"organization"  csv:"单位"`                                         // 单位，例如计算机学院，限15字
	Password     string `orm:"password"        json:"password" `                                                      // 密码
	RoleId       int    `orm:"role_id"         json:"roleId"`                                                         // 角色id
	Major        string `orm:"major"           json:"major"  csv:"专业"`                                                // 专业，限15字
}
