// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// LabCommentDao is the manager for logic model data accessing and custom defined data operations functions management.
type LabCommentDao struct {
	Table   string            // Table is the underlying table name of the DAO.
	Group   string            // Group is the database configuration group name of current DAO.
	Columns LabCommentColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// LabCommentColumns defines and stores column names for table lab_comment.
type LabCommentColumns struct {
	LabCommentId string // 主键
	LabId        string // 实验id
	CommentText  string // 评论内容，限120字
	Pid          string // 父评论id，主评时为0
	UserId       string // 发评论的用户id
	CreatedAt    string // 创建时间
	UpdatedAt    string // 更新时间
}

//  labCommentColumns holds the columns for table lab_comment.
var labCommentColumns = LabCommentColumns{
	LabCommentId: "lab_comment_id",
	LabId:        "lab_id",
	CommentText:  "comment_text",
	Pid:          "pid",
	UserId:       "user_id",
	CreatedAt:    "created_at",
	UpdatedAt:    "updated_at",
}

// NewLabCommentDao creates and returns a new DAO object for table data access.
func NewLabCommentDao() *LabCommentDao {
	return &LabCommentDao{
		Group:   "default",
		Table:   "lab_comment",
		Columns: labCommentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *LabCommentDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *LabCommentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *LabCommentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
