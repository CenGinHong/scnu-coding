// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// SysUserDao is the manager for logic model data accessing and custom defined data operations functions management.
type SysUserDao struct {
	Table   string         // Table is the underlying table name of the DAO.
	Group   string         // Group is the database configuration group name of current DAO.
	Columns SysUserColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// SysUserColumns defines and stores column names for table sys_user.
type SysUserColumns struct {
	UserId       string // 主键
	Email        string // 邮箱
	UserNum      string // 学号/职工号
	Username     string // 真实姓名
	Grade        string // 年级
	School       string // 学院
	Password     string // 密码
	AvatarImg    string // 头像url
	Gender       string // 性别
	Major        string // 专业
	Organization string // 单位，例如软件工程4班
	RoleId       string // 角色id
	UpdatedAt    string // 修改时间
	CreatedAt    string // 创建时间
}

//  sysUserColumns holds the columns for table sys_user.
var sysUserColumns = SysUserColumns{
	UserId:       "user_id",
	Email:        "email",
	UserNum:      "user_num",
	Username:     "username",
	Grade:        "grade",
	School:       "school",
	Password:     "password",
	AvatarImg:    "avatar_img",
	Gender:       "gender",
	Major:        "major",
	Organization: "organization",
	RoleId:       "role_id",
	UpdatedAt:    "updated_at",
	CreatedAt:    "created_at",
}

// NewSysUserDao creates and returns a new DAO object for table data access.
func NewSysUserDao() *SysUserDao {
	return &SysUserDao{
		Group:   "default",
		Table:   "sys_user",
		Columns: sysUserColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysUserDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysUserDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysUserDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
