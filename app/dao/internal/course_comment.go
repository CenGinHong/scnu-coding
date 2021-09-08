// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// CourseCommentDao is the manager for logic model data accessing and custom defined data operations functions management.
type CourseCommentDao struct {
	Table   string          // Table is the underlying table name of the DAO.
	Group   string          // Group is the database configuration group name of current DAO.
	Columns CourseCommentColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// CourseCommentColumns defines and stores column names for table course_comment.
type CourseCommentColumns struct {
	CommentId    string // 主键                  
    CourseId     string // 实验内容              
    CommentText  string // 评论内容，限120字     
    Pid          string // 父评论id，主评时为空  
    UserId       string // 发评论的用户id        
    CreatedAt    string // 创建时间              
    UpdatedAt    string // 更新时间
}

//  courseCommentColumns holds the columns for table course_comment.
var courseCommentColumns = CourseCommentColumns{
	CommentId:   "comment_id",    
            CourseId:    "course_id",     
            CommentText: "comment_text",  
            Pid:         "pid",           
            UserId:      "user_id",       
            CreatedAt:   "created_at",    
            UpdatedAt:   "updated_at",
}

// NewCourseCommentDao creates and returns a new DAO object for table data access.
func NewCourseCommentDao() *CourseCommentDao {
	return &CourseCommentDao{
		Group:   "default",
		Table:   "course_comment",
		Columns: courseCommentColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CourseCommentDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CourseCommentDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CourseCommentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}