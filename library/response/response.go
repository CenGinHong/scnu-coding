package response

// @Author: 陈健航
// @Date: 2020/11/1 22:55
// @Description:

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
)

// JsonResponse 数据返回通用JSON数据结构
type JsonResponse struct {
	Code    int         `json:"code"`    // 错误码(0成功，其他错误)
	Message string      `json:"message"` // 提示信息
	Data    interface{} `json:"result"`  // 返回数据(业务接口定义具体数据结构)
}

// PageInfo 分页参数
type PageInfo struct {
	Total   int `json:"total"`   // 总记录数
	Current int `json:"current"` // 当前页码
}

// PageResp 页面/**/返回集
type PageResp struct {
	Records interface{}              `json:"records"` // 业务数据
	Total   int                      `json:"total"`   // 业务总数
	Filter  map[string][]*FilterType `json:"filter"`  // 可筛选项
}

type FilterType struct {
	Text     string        `json:"text"`
	Value    string        `json:"value"`
	Children *[]FilterType `json:"children"`
}

// GetPageResp 返回分页结果集
func GetPageResp(records interface{}, total int, filter map[string][]*FilterType) (resp *PageResp) {
	recordsData := gconv.Interfaces(records)
	resp = &PageResp{
		Records: recordsData,
		Total:   total,
		Filter:  filter,
	}
	return resp
}

// Succ 成功返回结果集
// @params r
// @params data
// @date 2021-01-04 22:16:50
func Succ(r *ghttp.Request, data ...interface{}) {
	var respData = interface{}(nil)
	if len(data) > 0 {
		respData = data[0]
	}
	_ = r.Response.WriteJson(&JsonResponse{
		Code:    gerror.CodeOk,
		Message: "执行成功",
		Data:    respData,
	})
}

// Exit 发生错误返回结果集
// @params r
// @params error
// @date 2021-01-04 22:17:08
func Exit(r *ghttp.Request, err error) {
	//打印错误日志
	g.Log().Errorf("[url:%s][err:%s]",
		r.URL.Path, err.Error())
	// 封装错误信息,返回给前端
	_ = r.Response.WriteJson(
		JsonResponse{
			Code:    -1,
			Message: err.Error(),
			Data:    nil,
		},
	)
	r.Exit()
}
