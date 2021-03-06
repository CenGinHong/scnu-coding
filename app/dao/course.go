// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"scnu-coding/app/dao/internal"
)

// courseDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type courseDao struct {
	*internal.CourseDao
}

var (
	// Course is globally public accessible object for table course operations.
	Course courseDao
)

func init() {
	Course = courseDao{
		internal.NewCourseDao(),
	}
}

// Fill with you ideas below.
