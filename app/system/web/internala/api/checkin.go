package api

import (
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/system/web/internala/service"
	"scnu-coding/library/response"
)

var Checkin = checkinAPI{}

type checkinAPI struct{}

func (a checkinAPI) ListCheckinRecordByCourseId(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Checkin.ListCheckinRecordByCourseId(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

// StartCheckIn 教师开启签到
// @receiver a *checkinAPI
// @param r *ghttp.Request
// @date 2021-08-04 16:28:19
func (a *checkinAPI) StartCheckIn(r *ghttp.Request) {
	var req *define.StartCheckInReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	id, err := service.Checkin.StartCheckin(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{
		"id": id,
	})
}

// GetCheckinStatus 学生获取签到状态
// @receiver a *checkinAPI
// @param r *ghttp.Request
// @date 2021-08-04 16:48:21
func (a *checkinAPI) GetCheckinStatus(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Checkin.GetCheckinStatus(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

// StuListCheckinRecords 学生列表签到记录
// @receiver a *checkinAPI
// @param r *ghttp.Request
// @date 2021-08-04 16:49:02
func (a *checkinAPI) StuListCheckinRecords(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	resp, err := service.Checkin.StuListCheckinRecords(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

// ListCheckinDetailByCheckInRecordId 教师列表课程签到详情
// @receiver a *checkinAPI
// @param r *ghttp.Request
// @date 2021-08-04 16:49:10
func (a *checkinAPI) ListCheckinDetailByCheckInRecordId(r *ghttp.Request) {
	checkInRecordId := r.GetInt("checkinRecordId")
	resp, err := service.Checkin.ListCheckinDetailByCheckInRecordId(r.Context(), checkInRecordId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

// UpdateCheckinDetail 更新签到详情
// @receiver a *checkinAPI
// @param r *ghttp.Request
// @date 2021-08-04 16:49:57
func (a *checkinAPI) UpdateCheckinDetail(r *ghttp.Request) {
	var req *define.UpdateCheckinDetailReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.Checkin.UpdateCheckinDetail(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

// CheckIn 学生进行签到
// @receiver a *checkinAPI
// @param r *ghttp.Request
// @date 2021-08-04 16:50:12
func (a *checkinAPI) CheckIn(r *ghttp.Request) {
	var req *define.StudentCheckinReq
	if err := r.Parse(req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.Checkin.CheckIn(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

// DeleteCheckinRecord 删除签到记录
// @receiver a *checkinAPI
// @param r *ghttp.Request
// @date 2021-08-04 16:53:54
func (a *checkinAPI) DeleteCheckinRecord(r *ghttp.Request) {
	id := r.GetInt("id")
	if err := service.Checkin.DeleteCheckinRecord(r.Context(), id); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (a *checkinAPI) ExportCheckinRecord(r *ghttp.Request) {
	courseId := r.GetInt("courseId")
	file, err := service.Checkin.ExportCheckinToExcel(r.Context(), courseId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	r.Response.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	r.Response.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf8")
	r.Response.Header().Add("Content-Disposition", "attachment;filename="+gurl.Encode("签到.xlsx"))
	r.Response.WriteExit(file)
}
