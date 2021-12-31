package utils

import (
	"context"
	"errors"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/crypto/bcrypt"
	"scnu-coding/app/dao"
	"scnu-coding/app/model"
	"scnu-coding/library/response"
)

// @Author: 陈健航
// @Date: 2021/5/1 16:44
// @Description:

// GfToken GfToken
var GfToken = newGfToken()

func newGfToken() (gfToken gtoken.GfToken) {
	authExcludePaths := g.SliceStr{
		"/web/user/nickname",
		"/web/user/signup/stu",
		"/web/user/email",
		"/web/user/password",
		"/web/user/verificationCode",
		"/web/user/test/*",
		"/web/test/*",
		"/web/ide/connect",
	}
	gfToken = gtoken.GfToken{
		CacheMode:        2,
		LoginPath:        "/login",
		LogoutPath:       "/logout",
		LoginBeforeFunc:  LoginBeforeFunc,
		LoginAfterFunc:   LoginAfterFunc,
		LogoutAfterFunc:  LogoutAfterFunc,
		AuthExcludePaths: authExcludePaths,
	}
	if g.Cfg().GetBool("server.IsMultiple") {
		gfToken.CacheMode = 2
	}
	return gfToken
}

// AuthAfterFunc 身份认证操作的后续,主要是鉴权
// @params r
// @params respData
// @date 2021-01-04 22:13:39
//func AuthAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
// 存在令牌
//if respData.Success() {
//	// 鉴权
//	Id := respData.GetString("userKey")
//	if ok, err := authenticate(Id, r.URL.Path, r.Method); err != nil {
//		response.Exit(r, code.OtherError)
//	} else if !ok {
//		// 权限不足
//		response.Exit(r, code.PermissionError)
//	}
//	// 鉴权成功
//	ctxUser := &model.ContextUser{}
//	if err := g.Model("sys_user").InnerJoin("sys_re_user_role").
//		InnerJoin("sys_role").Where("sys_user.user_id", Id).
//		Where("sys_re_user_role.user_id = sys_user.user_id").Where("sys_re_user_role.user_id = sys_role.role_id").
//		Fields(&ctxUser).Cache(1 * time.Minute).Scan(&ctxUser); err != nil {
//		return
//	}
//	service.Context.Get(r.Context()).User = ctxUser
//	r.Middleware.Next()
//	//不存在令牌
//} else {
//	var params map[string]interface{}
//	if r.Method == "GET" {
//		params = r.GetMap()
//	} else if r.Method == "POST" {
//		params = r.GetMap()
//	} else {
//		response.Exit(r, code.OtherError)
//		return
//	}
//	no := gconv.String(gtime.TimestampMilli())
//	g.Log().Info("[AUTH_%s][url:%s][params:%s][data:%s]",
//		no, r.URL.Path, params, respData.Json())
//	response.Exit(r, code.AuthError)
//}
//}

//// authenticate 鉴权方法
//// @params userId
//// @params url
//// @params method
//// @return bool
//// @return error
//// @date 2021-01-04 21:58:46
//func authenticate(Id string, url string, method string) (ok bool, err error) {
//	// 定义结构体
//	type API struct {
//		API    string
//		Method string
//	}
//	// 创建SQL结果集
//	apis := make([]API, 0)
//	if err = g.Model("sys_api").InnerJoin("sys_re_api_role").InnerJoin("sys_re_user_role").
//		InnerJoin("sys_user").InnerJoin("sys_role").Cache(5*time.Minute).
//		Where("sys_user.user_id =", Id).And("sys_user.user_id = sys_re_user_role.user_id").
//		Where("sys_user_role.role_id = sys_re_api_role_id.role_id").Where("sys_re_api_role.api_id = sys_api.api_id").
//		Fields("api", "method").Scan(&apis); err != nil {
//		return
//	}
//
//	for _, v := range apis {
//		// 用正则匹配,在该用户角色的可访问接口里有无该API
//		if isMatch := gregex.IsMatchString(v.API, url) && v.Method == method; isMatch {
//			return
//		}
//	}
//	// 该用户角色没有权限访问该接口
//	return
//}

// LoginBeforeFunc 登录方法
func LoginBeforeFunc(r *ghttp.Request) (string, interface{}) {
	//var req *struct {
	//	UserNum  string `valid:"required|integer#登录名不能为空|登录名需要为学号"` // 登陆凭证，目前是email
	//	Password string `valid:"required|password#密码不能为空|密码不符合规则"`  // 密码
	//}
	var req *struct {
		UserNum  string `valid:"required|integer#登录名不能为空|登录名需要为学号"` // 登陆凭证，目前是email
		Password string `valid:"required#密码不能为空"`                   // 密码
	}
	// 转换成结构体
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return "", nil
	}
	// 在数据库查询用户是否存在(只查出密码）
	password, err := dao.SysUser.Ctx(context.TODO()).Where(dao.SysUser.Columns.UserNum, req.UserNum).Value(dao.SysUser.Columns.Password)
	if err != nil {
		response.Exit(r, err)
		return "", nil
	}
	// 不存在该用户
	if password.IsNil() {
		response.Exit(r, errors.New("账号或密码错误"))
		return "", nil
	}
	// 校验密码
	if err = bcrypt.CompareHashAndPassword(password.Bytes(), []byte(req.Password)); err != nil {
		response.Exit(r, errors.New("账号或密码错误"))
	}
	// 获取信息
	userInfo := &model.ContextUser{}
	if err = dao.SysUser.Ctx(context.TODO()).Where(dao.SysUser.Columns.UserNum, req.UserNum).Scan(&userInfo); err != nil {
		response.Exit(r, err)
		return "", nil
	}
	//校验成功
	return gconv.String(userInfo.UserId), userInfo
}

// LoginAfterFunc 重定义返回后结果集
// @params r
// @params respData
// @date 2021-01-04 22:14:51
func LoginAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	if respData.Success() {
		// 返回token
		response.Succ(r, g.Map{
			"token": respData.GetString("token"),
		})
	}
}

// LogoutAfterFunc 重定义退登结果集
// @params r
// @params respData
// @date 2021-01-04 22:14:32
func LogoutAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	if respData.Success() {
		response.Succ(r)
	} else {
		response.Exit(r, errors.New("退出登陆失败"))
	}
}
