package api

import (
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/system/web/internala/service"
	"scnu-coding/library/response"
)

// LabSummit service
var LabSummit = labSummitAPI{}

type labSummitAPI struct{}

func (l *labSummitAPI) UpdateFinishStat(r *ghttp.Request) {
	var req *define.UpdateLabFinishReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.LabSummit.UpdateFinishStat(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (l *labSummitAPI) GetReportContent(r *ghttp.Request) {
	var req *define.GetReportContentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	resp, err := service.LabSummit.GetReportContent(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, g.Map{
		"reportContent": resp,
	})
}

func (l *labSummitAPI) UpdateReportContent(r *ghttp.Request) {
	var req *define.UpdateReportContentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	if err := service.LabSummit.UpdateReport(r.Context(), req); err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, true)
}

func (l *labSummitAPI) ListLabSubmit(r *ghttp.Request) {
	labId := r.GetInt("labId")
	resp, err := service.LabSummit.ListLabSummit(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (l *labSummitAPI) ListLabSubmitId(r *ghttp.Request) {
	labId := r.GetInt("labId")
	resp, err := service.LabSummit.ListLabSummitId(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (l *labSummitAPI) UpdateScoreAndComment(r *ghttp.Request) {
	var req *define.UpdateLabSummitScoreAndCommentReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	err := service.LabSummit.UpdateScoreAndComment(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r)
}

func (l *labSummitAPI) ExportScore(r *ghttp.Request) {
	var req *define.ExportLabScoreReq
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return
	}
	file, err := service.LabSummit.ExportScore(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	r.Response.Header().Add("Access-Control-Expose-Headers", "Content-Disposition")
	r.Response.Header().Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=utf8")
	r.Response.Header().Add("Content-Disposition", "attachment;filename="+gurl.Encode("成绩.xlsx"))
	r.Response.WriteExit(file)
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
		return
	}
	resp, err := service.LabSummit.GetCodeData(r.Context(), req)
	if err != nil {
		response.Exit(r, err)
		return
	}
	response.Succ(r, resp)
}

func (l *labSummitAPI) PlagiarismCheck(r *ghttp.Request) {
	labId := r.GetInt("labId")
	resp, err := service.LabSummit.ExecPlagiarismCheckByMoss(r.Context(), labId)
	if err != nil {
		response.Exit(r, err)
	}
	response.Succ(r, resp)
}

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
