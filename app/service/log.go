package service

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/dao"
)

// @Author: 陈健航
// @Date: 2021/5/6 13:12
// @Description:

var Log = logService{}

type logService struct{}

func (receiver *logService) Log(r *ghttp.Request) {
	// 记录除了查询意外的日志
	if r.Method != "GET" {
		requestParams, _ := gjson.LoadJson(r.GetMapStrStr())
		responseData := r.Response.BufferString()
		ctxUser := Context.Get(r.Context()).User
		log := g.Map{
			dao.Log.Columns.UserId:     ctxUser.UserId,
			dao.Log.Columns.HttpMethod: r.Method,
			dao.Log.Columns.ReqParams:  requestParams,
			dao.Log.Columns.RespData:   responseData,
		}
		// 插入日志
		_, _ = dao.Log.Ctx(r.Context()).Data(log).Insert()
	}
}
