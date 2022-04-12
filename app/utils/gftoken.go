package utils

import (
	"context"
	"errors"
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
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
		CacheMode:        1,
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
func AuthAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	//存在令牌
	if respData.Success() {
		// 鉴权
		_ = respData.GetString("userKey")
		//if ok, err := authenticate(r.Context(), id, r.Router.Uri, r.Method); err != nil {
		//	response.Exit(r, gerror.NewCode(gcode.CodeNotAuthorized))
		//} else if !ok {
		//	// 权限不足
		//	response.Exit(r, gerror.NewCode(gcode.CodeNotAuthorized))
		//}
		r.Middleware.Next()
	} else {
		//不存在令牌
		var params map[string]interface{}
		params = r.GetMap()
		no := gconv.String(gtime.TimestampMilli())
		g.Log().Info("[AUTH_%s][url:%s][params:%s][data:%s]",
			no, r.URL.Path, params, respData.Json())
		response.Exit(r, gerror.NewCode(gcode.CodeNotAuthorized))
	}
}

//// authenticate 鉴权方法
//// @params userId
//// @params url
//// @params method
//// @return bool
//// @return error
//// @date 2021-01-04 21:58:46
func authenticate(ctx context.Context, id string, uri string, method string) (ok bool, err error) {
	// 定义结构体
	type API struct {
		Path   string
		Method string
		Allow  []int
	}
	// 创建SQL结果集
	apis := make([]API, 0)
	roleId, err := dao.SysUser.Ctx(ctx).WherePri(id).Value(dao.SysUser.Columns.RoleId)
	if err != nil {
		return false, err
	}
	contents := gfile.GetContentsWithCache("/var/www/scnu-coding/config/authority.json", 0)
	if err = gjson.DecodeTo(contents, &apis); err != nil {
		return false, err
	}
	isAllow := false
	for _, api := range apis {
		if api.Path == uri && api.Path == method {
			for _, i := range api.Allow {
				if i == roleId.Int() {
					isAllow = true
					break
				}
			}
		}
	}
	// 该用户角色没有权限访问该接口
	return isAllow, nil
}

// LoginBeforeFunc 登录方法
func LoginBeforeFunc(r *ghttp.Request) (string, interface{}) {
	var req *struct {
		UserNum  string `valid:"required|integer#登录名不能为空|登录名需要为学号"` // 登陆凭证，目前是学号
		Password string `valid:"required#密码不能为空"`                   // 密码
	}
	// 转换成结构体
	if err := r.Parse(&req); err != nil {
		response.Exit(r, err)
		return "", nil
	}
	// 在数据库查询用户是否存在(只查出密码）
	password, err := dao.SysUser.Ctx(r.Context()).Where(dao.SysUser.Columns.UserNum, req.UserNum).
		Value(dao.SysUser.Columns.Password)
	if err != nil {
		response.Exit(r, err)
		return "", nil
	}
	// 不存在该用户
	if password.IsNil() {
		response.Exit(r, gerror.NewCode(gcode.CodeNotAuthorized))
		return "", nil
	}
	// 校验密码
	if err = bcrypt.CompareHashAndPassword(password.Bytes(), []byte(req.Password)); err != nil {
		response.Exit(r, errors.New("账号或密码错误"))
	}
	// 获取信息
	userInfo := &model.ContextUser{}
	if err = dao.SysUser.Ctx(r.Context()).
		Where(dao.SysUser.Columns.UserNum, req.UserNum).
		Scan(&userInfo); err != nil {
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
