package service

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/model"
	"scnu-coding/app/utils"
)

// Middleware 中间件管理服务
var Middleware = serviceMiddleware{}

type serviceMiddleware struct{}

func (s *serviceMiddleware) Ctx(r *ghttp.Request) {
	// 初始化，务必最开始执行
	customCtx := &model.Context{}
	Context.Init(r, customCtx)
	if r.RequestURI != "/login" {
		data := utils.GfToken.GetTokenData(r)
		contextUser := &model.ContextUser{}
		_ = gconv.Struct(data.Get("data"), &contextUser)
		customCtx.User = contextUser
	}
	// 解析分页数据
	pageInfo := &model.ContextPageInfo{}
	_ = r.ParseQuery(&pageInfo)
	if pageInfo.PageSize == 0 {
		pageInfo.PageSize = 10
	}
	if pageInfo.Current == 0 {
		pageInfo.Current = 1
	}
	pageInfo.SortField = gstr.CaseSnake(pageInfo.SortField)
	pageInfo.SortOrder = gstr.Replace(pageInfo.SortOrder, "end", "")
	pageInfo.ParseFilterFields = make(map[string][]string, 0)
	// 解析转换
	filterFields := r.Get("filterFields")
	temp := gconv.MapStrStr(filterFields)
	for k, v := range temp {
		k = gstr.CaseSnake(k)
		//这里的temp1是{
		//	1:计算机
		//	2：通信这样的结构
		//}
		temp1 := gconv.MapStrStr(v)
		for _, v1 := range temp1 {
			pageInfo.ParseFilterFields[k] = append(pageInfo.ParseFilterFields[k], v1)
		}
	}
	customCtx.PageInfo = pageInfo
	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// CORS 允许接口跨域请求
func (s *serviceMiddleware) CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
