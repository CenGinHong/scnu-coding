// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// CourseAnnouncementDao is the manager for logic model data accessing and custom defined data operations functions management.
type CourseAnnouncementDao struct {
	Table   string                    // Table is the underlying table name of the DAO.
	Group   string                    // Group is the database configuration group name of current DAO.
	Columns CourseAnnouncementColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// CourseAnnouncementColumns defines and stores column names for table course_announcement.
type CourseAnnouncementColumns struct {
	CourseAnnouncementId string // id
	Title                string // 标题
	CourseId             string // 课程id
	Content              string // 公告内容，限2000字
	AttachmentSrc        string // 文件url
	CreatedAt            string // 创建时间
	UpdatedAt            string // 修改时间
}

//  courseAnnouncementColumns holds the columns for table course_announcement.
var courseAnnouncementColumns = CourseAnnouncementColumns{
	CourseAnnouncementId: "course_announcement_id",
	Title:                "title",
	CourseId:             "course_id",
	Content:              "content",
	AttachmentSrc:        "attachment_src",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

// NewCourseAnnouncementDao creates and returns a new DAO object for table data access.
func NewCourseAnnouncementDao() *CourseAnnouncementDao {
	return &CourseAnnouncementDao{
		Group:   "default",
		Table:   "course_announcement",
		Columns: courseAnnouncementColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CourseAnnouncementDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CourseAnnouncementDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CourseAnnouncementDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
