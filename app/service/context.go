package service

import (
	"context"
	"github.com/gogf/gf/net/ghttp"
	"scnu-coding/app/model"
)

// Context 上下文管理服务
var Context = contextService{
	contextKey: "contextKey", // 上下文变量存储键名
}

type contextService struct {
	contextKey string
}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
// @receiver s
// @params r
// @params customCtx
// @date 2021-05-01 12:20:28
func (s *contextService) Init(r *ghttp.Request, customCtx *model.Context) {
	r.SetCtxVar(s.contextKey, customCtx)
}

// Get 获得上下文变量，如果没有设置，那么返回nil
// @receiver s
// @params ctx
// @return *model.Context
// @date 2021-05-01 12:20:33
func (s *contextService) Get(ctx context.Context) *model.Context {
	value := ctx.Value(s.contextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*model.Context); ok {
		return localCtx
	}
	return nil
}
