// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// CheckinRecordDao is the manager for logic model data accessing and custom defined data operations functions management.
type CheckinRecordDao struct {
	Table   string          // Table is the underlying table name of the DAO.
	Group   string          // Group is the database configuration group name of current DAO.
	Columns CheckinRecordColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// CheckinRecordColumns defines and stores column names for table checkin_record.
type CheckinRecordColumns struct {
	CheckinRecordId  string // id                              
    CheckinName      string // 签到名称，例如2021年2月5日签到  
    CheckinKey       string // 签到密钥                        
    TotalDuration    string // 限时时间                        
    CourseId         string // 课程id                          
    CreatedAt        string // 创建时间                        
    UpdatedAt        string // 更新时间
}

//  checkinRecordColumns holds the columns for table checkin_record.
var checkinRecordColumns = CheckinRecordColumns{
	CheckinRecordId: "checkin_record_id",  
            CheckinName:     "checkin_name",       
            CheckinKey:      "checkin_key",        
            TotalDuration:   "total_duration",     
            CourseId:        "course_id",          
            CreatedAt:       "created_at",         
            UpdatedAt:       "updated_at",
}

// NewCheckinRecordDao creates and returns a new DAO object for table data access.
func NewCheckinRecordDao() *CheckinRecordDao {
	return &CheckinRecordDao{
		Group:   "default",
		Table:   "checkin_record",
		Columns: checkinRecordColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *CheckinRecordDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *CheckinRecordDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *CheckinRecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}