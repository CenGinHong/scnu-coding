package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/system/web/internala/service"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/5/24 23:35
// @Description:

var Course = courseAPI{}

type courseAPI struct{}

// GetOne 获取某课程的详情
// @Description:
// @receiver c
// @param r
// @date 2021-10-17 13:13:34
func (c *courseAPI) GetOne(r *ghttp.Request) {
	courseId := r.GetInt("id")
	resp, err := service.Course.GetCourseDetail(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

// ListCourseStudentOverview 列表课程学生的信息
// @Description:
// @receiver c
// @param r
// @date 2021-10-17 13:14:01
func (c *courseAPI) ListCourseStudentOverview(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Course.ListCourseStudentOverview(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

// ListCourseEnroll
// @Description:
// @receiver c
// @param r
// @date 2021-10-17 13:14:34
func (c *courseAPI) ListCourseEnroll(r *ghttp.Request) {
	resp, err := service.Course.ListCourseEnroll(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (c *courseAPI) IsCourseEnroll(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	isEnrollCourse, err := service.Course.IsEnrollCourse(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{
		"isEnroll": isEnrollCourse,
	})
}

func (c *courseAPI) IsOpenByTeacherId(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	isOpen, err := service.Course.IsOpenByTeacherId(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{
		"isOpen": isOpen,
	})
}

func (c *courseAPI) ImportStudent2Class(r *ghttp.Request) {
	var req *define.ImportStudent2Class
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	resp, err := service.Course.ImportStudent2Class(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

// InsertCourse 创建课程
// @receiver c *courseAPI
// @param r *ghttp.Request
// @date 2021-09-16 20:15:45
func (c *courseAPI) InsertCourse(r *ghttp.Request) {
	var req *define.InsertCourseReq
	if err := r.Parse(&req); err != nil {
		return
	}
	id, err := service.Course.InsertCourse(r.Context(), req)
	if err != nil {
		return
	}
	response.Succ(r, g.Map{
		"id": id,
	})
}

// UpdateCourse 更新课程
// @receiver receiver
// @params r
// @date 2021-05-25 00:01:19
func (c *courseAPI) UpdateCourse(r *ghttp.Request) {
	var req *define.UpdateCourseReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	//保存
	if err := service.Course.UpdateCourse(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, true)
}

func (c *courseAPI) Delete(r *ghttp.Request) {
	courseId := r.GetInt("id")
	err := service.Course.DeleteCourse(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

// ListCourseByTeacherId 列出教师开展的课程
// @receiver receiver
// @params r
// @date 2021-05-25 00:01:28
func (c *courseAPI) ListCourseByTeacherId(r *ghttp.Request) {
	resp, err := service.Course.ListCourseByTeacherId(r.Context())
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (c *courseAPI) JoinClass(r *ghttp.Request) {
	var req *define.JoinClassReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.Course.JoinClass(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

//// ListCourseEnroll 学生列出自己加入的课程
//// @receiver receiver
//// @params r
//// @date 2021-05-25 00:01:56
//func (c *courseAPI) ListCourseEnroll(r *ghttp.Request) {
//	resp, err := service.Course.ListCourseEnroll(r.Context())
//	if err != nil {
//		response.Exit(r, err)
//return
//	}
//	response.Succ(r, resp)
//}
//
//// ListStuByCourseId 列出所有选出某门课的学生
//// @receiver receiver
//// @params r
//// @date 2021-05-25 00:03:14
//func (c *courseAPI) ListStuByCourseId(r *ghttp.Request) {
//	CourseId := r.GetInt("courseId")
//	resp, err := service.Course.ListStudentByCourseId(r.Context(), CourseId)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, resp)
//}
//
//// Enroll 加入课程
//// @receiver receiver
//// @params r
//// @date 2021-05-25 00:04:12
//func (c *courseAPI) Enroll(r *ghttp.Request) {
//	var req *define.EnrollCourseReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	if err := service.Course.Enroll(r.Context(), req); err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, true)
//}
//
//// QuitCourse 删除选课记录
//// @receiver receiver
//// @params r
//// @date 2021-05-25 00:06:56
//func (c *courseAPI) QuitCourse(r *ghttp.Request) {
//	var req *define.DropCourseReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	if err := service.Course.DeleteEnrollRecord(r.Context(), req); err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, true)
//}
//
//func (c *courseAPI) DeleteCourse(r *ghttp.Request) {
//	courseId := r.GetInt("courseId")
//	if err := service.Course.Delete(r.Context(), courseId); err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, true)
//}

// ListCourseByCourseName 搜索课程
// @receiver receiver
// @params r
// @date 2021-05-25 00:09:45
func (c *courseAPI) ListCourseByCourseName(r *ghttp.Request) {
	courseName := r.GetString("courseName")
	resp, err := service.Course.ListCourseByCourseName(r.Context(), courseName)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (c *courseAPI) ListOneScore(r *ghttp.Request) {
	var req *define.ListOneStudentScoreReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	resp, err := service.Course.ListOneStudentScore(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

//
//func (c *courseAPI) ListAllCourse(r *ghttp.Request) {
//	resp, err := service.Course.ListAllCourse(r.Context())
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, resp)
//}
