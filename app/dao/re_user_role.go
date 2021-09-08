// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"scnu-coding/app/dao/internal"
)

// reUserRoleDao is the manager for logic model data accessing and custom defined data operations functions management. 
// You can define custom methods on it to extend its functionality as you wish.
type reUserRoleDao struct {
	*internal.ReUserRoleDao
}

var (
	// ReUserRole is globally public accessible object for table re_user_role operations.
	ReUserRole reUserRoleDao
)

func init() {
	ReUserRole = reUserRoleDao{
		internal.NewReUserRoleDao(),
	}
}

// Fill with you ideas below.