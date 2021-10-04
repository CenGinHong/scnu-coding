package api

import (
	"fmt"
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

func (l *labSummitAPI) ListLabSubmitId(r *ghttp.Request) {
	labId := r.GetInt("labId")
	resp, err := service.LabSummit.ListLabSummitId(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

func (l *labSummitAPI) UpdateScoreAndComment(r *ghttp.Request) {
	var req *define.UpdateLabSummitScoreAndCommentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	err := service.LabSummit.UpdateScoreAndComment(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r)
}

func (l *labSummitAPI) ExportScore(r *ghttp.Request) {
	var req *define.ExportLabScoreReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	file, err := service.LabSummit.ExportScore(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	r.Response.Header().Set("Pragma", "No-cache")
	r.Response.Header().Set("Cache-Control", "No-cache")
	r.Response.Header().Set("Expires", "0")
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "成绩表表"))
	r.Response.Header().Set("Content-Type", "text/csv")
	r.Response.Write(file)
	r.Response.Flush()
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

func (l *labSummitAPI) GetCode(r *ghttp.Request) {
	var req *define.ReadCodeDataReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
	}
	resp, err := service.LabSummit.GetCodeData(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

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
