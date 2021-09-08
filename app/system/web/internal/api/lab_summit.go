package api

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/system/web/internal/service"
	"scnu-coding/library/response"
)

// LabSummit service
var LabSummit = labSummitAPI{}

type labSummitAPI struct{}

func (l *labSummitAPI) UpdateFinishStat(r *ghttp.Request) {
	var req *define.UpdateLabFinishReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	if err := service.LabSummit.UpdateFinishStat(r.Context(), req); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r)
}

func (l *labSummitAPI) GetReportContent(r *ghttp.Request) {
	var req *define.GetReportContentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	resp, err := service.LabSummit.GetReportContent(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, g.Map{
		"reportContent": resp,
	})
}

func (l *labSummitAPI) UpdateReportContent(r *ghttp.Request) {
	var req *define.UpdateReportContentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	if err := service.LabSummit.UpdateReport(r.Context(), req); err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, true)
}

func (l *labSummitAPI) ListLabSubmit(r *ghttp.Request) {
	labId := r.GetInt("labId")
	resp, err := service.LabSummit.ListLabSummit(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

//func (l *labSummitAPI) InsertCodeFinish(r *ghttp.Request) {
//	var req *define.UpdateLabFinishReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	if err := service.LabSummit.UpdateFinishStat(r.Context(), req); err != nil {
//		response.Exit(r, err)
//	}
//}
//
//func (l *labSummitAPI) ListLabSummit(r *ghttp.Request) {
//	labId := r.GetInt("labId")
//	resp, err := service.LabSummit.ListLabSummit(r.Context(), labId)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, resp)
//}
//
//func (l *labSummitAPI) CollectCompilerErrorLog(r *ghttp.Request) {
//	labId := r.GetInt("labId")
//	log, err := service.LabSummit.CollectCompilerErrorLog(labId)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, g.Map{
//		"compiler_error_log": log,
//	})
//}
//
//func (l *labSummitAPI) readCodeData(r *ghttp.Request) {
//	var req *define.CheckCodeQuickReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	resp, err := service.LabSummit.readCodeData(r.Context(), req)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, resp)
//}

//func (receiver *labSummitAPI) PlagiarismCheck(r *ghttp.Request) {
//	labId := r.GetInt("labId")
//	resp, err := service.LabSummit.PlagiarismCheck(labId)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, resp)
//}

//func (receiver *labSummitAPI) GetReportUrl(r *ghttp.Request) {
//	var req *model.GetReportUrlReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	req.StuID = r.GetCtxVar(dao.SysUser.Columns.StuID).Int()
//	url, err := service.LabSummit.GetReportUrl(req)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, g.Map{"url": url})
//}

//func (l labSummitAPI) UpdateScoreAndComment(r *ghttp.Request) {
//	var req *define.UpdateLabSummitScoreReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	if err := service.LabSummit.UpdateScoreAndComment(req); err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r)
//}
//
//func (l labSummitAPI) UpdateComment(r *ghttp.Request) {
//	var req *define.UpdateLabSummitCommentReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	if err := service.LabSummit.UpdateComment(req); err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r)
//}

//func (receiver *labSummitAPI) GetCommentByStuId(r *ghttp.Request) {
//	var req *model.GetCommentByStuIdReq
//	if err := r.Parse(&req); err != nil {
//		response.Exit(r, err)
//	}
//	comment, err := service.LabSummit.GetCommentByStuId(req)
//	if err != nil {
//		response.Exit(r, err)
//	}
//	response.Succ(r, g.Map{
//		"comment": comment,
//	})
//}
