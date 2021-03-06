// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"scnu-coding/app/dao/internal"
)

// logDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type logDao struct {
	*internal.LogDao
}

var (
	// Log is globally public accessible object for table log operations.
	Log logDao
)

func init() {
	Log = logDao{
		internal.NewLogDao(),
	}
}

// Fill with you ideas below.
