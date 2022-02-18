package model

// Context 请求上下文结构
type Context struct {
	PageInfo *ContextPageInfo // 文件
	User     *ContextUser     // 上下文用户信息
}

// ContextUser 请求上下文中的用户信息
type ContextUser struct {
	UserId int // 用户ID
	RoleId int // 用户角色
}

// ContextPageInfo 请求上下文中的页面
type ContextPageInfo struct {
	Current           int                 // 当前页码
	PageSize          int                 // 页面大小
	SortField         string              // 排序域名
	SortOrder         string              // 排序顺序
	ParseFilterFields map[string][]string //解析后的数据，业务中使用
}
