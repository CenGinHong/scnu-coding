package service

import (
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/model"
	"scnu-coding/app/utils"
	"scnu-coding/library/response"
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
	var pageInfo *model.ContextPageInfo
	_ = r.ParseQuery(&pageInfo)
	if pageInfo != nil {
		pageInfo.SortField = gstr.CaseSnake(pageInfo.SortField)
		pageInfo.SortOrder = gstr.Replace(pageInfo.SortOrder, "end", "")
		pageInfo.ParseFilterFields = make(map[string][]string, 0)
		tempFilterFields := make(map[string][]string, 0)
		// 解析转换
		if err := gconv.Struct(r.Get("filterFields"), &tempFilterFields); err != nil {
			response.Exit(r, err)
		}
		// 转成驼峰方便数据库查询
		for key, value := range tempFilterFields {
			pageInfo.ParseFilterFields[gstr.CaseSnake(key)] = value
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
