package service

// @Author: 陈健航
// @Date: 2020/12/31 0:10
// @Description:

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"golang.org/x/crypto/bcrypt"
	"scnu-coding/app/dao"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internal/define"
)

const (
	Teacher = 1 + iota
	Student
)

var SysUser = userService{}

type userService struct{}

func (u *userService) IsEmailUsed(ctx context.Context, email string) (isUsed bool, err error) {
	ctxUser := service.Context.Get(ctx).User
	count, err := dao.SysUser.Ctx(ctx).WhereNot(dao.SysUser.Columns.Email, email).WherePri(ctxUser.UserId).Count()
	if err != nil {
		return
	}
	isUsed = count > 0
	return isUsed, nil
}

func (u *userService) IsUserNumUsed(ctx context.Context, userNum string) (isUsed bool, err error) {
	ctxUser := service.Context.Get(ctx).User
	count, err := dao.SysUser.Ctx(ctx).WhereNot(dao.SysUser.Columns.UserId, ctxUser.UserId).Where(dao.SysUser.Columns.UserNum, userNum).Count()
	if err != nil {
		return false, err
	}
	isUsed = count > 0
	return isUsed, nil
}

func (u userService) GetUserInfo(ctx context.Context) (resp *define.GetUserInfoResp, err error) {
	ctxUser := service.Context.Get(ctx).User
	resp = &define.GetUserInfoResp{}
	if err = dao.SysUser.Ctx(ctx).WherePri(ctxUser.UserId).Scan(&resp); err != nil {
		return nil, err
	}
	return resp, nil

}

// StudentRegister 学生注册
// @receiver s
// @params req
// @return error
// @date 2021-01-09 00:15:22
//func (u *userService) StudentRegister(_ context.Context, req *define.RegisterReq) (err error) {
//	// 校验验证码
//	if err = u.checkVerCode(req.Email, req.VerCode); err != nil {
//		return err
//	}
//	// 密码加密
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//	//存入加密后的密码
//	data := g.Map{
//		dao.SysUser.Columns.Email:    req.Email,
//		dao.SysUser.Columns.Password: hashedPassword,
//		dao.SysUser.Columns.Username: req.Username,
//	}
//	// 保存
//	lastInsertId, err := dao.SysUser.InsertAndGetId(data)
//	if err != nil {
//		return err
//	}
//	// 赋予权限
//	if _, err = dao.ReUserRole.InsertLab(g.Map{
//		dao.ReUserRole.Columns.UserId: lastInsertId,
//		dao.ReUserRole.Columns.RoleId: Student,
//	}); err != nil {
//		return err
//	}
//	return nil
//}

// TeacherRegister 教师签发账户注册
// @receiver s
// @params req
// @return error
// @date 2021-01-09 00:15:22
//func (u *userService) TeacherRegister(req *define.RegisterReq) error {
//	// 不需要验证码
//	// 密码加密
//	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//	//存入加密后的密码
//	req.Password = string(hashPassword)
//	if _, err = dao.SysUser.Data(g.Map{
//		dao.SysUser.Columns.Email:    req.Email,
//		dao.SysUser.Columns.Password: hashPassword,
//		dao.SysUser.Columns.Username: req.Username,
//	}).InsertLab(); err != nil {
//		return err
//	}
//	return nil
//}

// Update 修改个人资料
// @receiver s
// @params req
// @return error
// @date 2021-01-10 00:07:55
func (u *userService) Update(ctx context.Context, req *define.UpdateUserInfoReq) (err error) {
	ctxUser := service.Context.Get(ctx).User
	oldEmail, err := dao.SysUser.Ctx(ctx).WherePri(ctxUser.UserId).Value(dao.SysUser.Columns.Email)
	if err != nil {
		return err
	}
	// 修改了邮箱,检查邮箱是否验证成功
	if req.Email != oldEmail.String() {
		if err = service.Common.CheckVerCode(req.Email, req.VerCode); err != nil {
			return err
		}
	}
	// 修改了密码，检查密码是否
	if req.Password != "" {
		password, err := dao.SysUser.Ctx(ctx).WherePri(ctxUser.UserId).Value(dao.SysUser.Columns.Password)
		if err != nil {
			return err
		}
		if err = bcrypt.CompareHashAndPassword(password.Bytes(), []byte(req.OldPassword)); err != nil {
			return gerror.NewCode(-1, "密码验证错误")
		}
	}
	if _, err = dao.SysUser.Ctx(ctx).WherePri(ctxUser.UserId).Data(req).Update(); err != nil {
		return err
	}
	return nil
}

// ResetPassword 重置密码
// @receiver s
// @params req
// @return error
// @date 2021-01-10 00:07:38
//func (u *userService) ResetPassword(_ context.Context, req *define.ResetPasswordReq) error {
//	// 检查验证码
//	if err := u.checkVerCode(req.Email, req.VerCode); err != nil {
//		return err
//	}
//	// 密码加密
//	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
//	if err != nil {
//		return err
//	}
//	// 存入加密后的密码
//	if _, err = dao.SysUser.Where(dao.SysUser.Columns.Email, req.Email).
//		Update(dao.SysUser.Columns.Password, hashedPassword); err != nil {
//		return err
//	}
//	return nil
//}

// DeleteUser 注销用户
// @receiver s
// @params req
// @return error
// @date 2021-01-14 11:27:40
//func (u *userService) DeleteUser(ctx context.Context, req *define.DeleteUserReq) (err error) {
//	ctxUser := service.Context.Get(ctx).User
//	userPassword, err := dao.SysUser.WherePri(ctxUser.UserId).Value(dao.SysUser.Columns.Password)
//	if err != nil {
//		return err
//	}
//	// 是否存在该用户
//	if !userPassword.IsNil() {
//		return code.UserNotExistError
//	}
//	// 校验密码
//	if err = bcrypt.CompareHashAndPassword(userPassword.Bytes(), []byte(req.Password)); err != nil {
//		return code.PasswordError
//	}
//	// 执行删除
//	if _, err = dao.SysUser.Delete(ctxUser.UserId); err != nil {
//		return err
//	}
//	return nil
//}

//func (receiver *userService) GetUserInfoByToken(id int) (resp *model.SysUserResp, err error) {
//	resp = &model.SysUserResp{}
//	if err = dao.SysUser.WherePri(id).FieldsEx(dao.SysUser.Columns.DeletedAt).Scan(&resp); err != nil {
//		return nil, err
//	}
//	if resp.Role, err = dao.SysUser.GetRoleById(id); err != nil {
//		return nil, err
//	}
//	return resp, err
//}

// ListCodingTimeByUserId
// @receiver receiver
// @params ctx
// @params year 年份
// @return resp
// @return err
// @date 2021-05-02 17:01:39
func (u *userService) ListCodingTimeByUserId(ctx context.Context, req *define.ListCodingTimeByUserIdReq) (resp []*define.CodingTimeRecord, err error) {
	resp = make([]*define.CodingTimeRecord, 0)
	ctxUser := service.Context.Get(ctx).User
	if req.UserId == 0 {
		req.UserId = ctxUser.UserId
	}
	d := dao.CodingTime.Ctx(ctx).Fields("SMU(duration)", "Date_Format(created_at,'%Y-%m-%d') as created_at").
		Group("Date_Format(created_at,'%Y-%m-%d')").Where(dao.CodingTime.Columns.UserId, req.UserId)
	// 限定日期范围
	if req.BeginDate != "" {
		d = d.WhereGTE(dao.CodingTime.Columns.CreatedAt, req.BeginDate)
	}
	if req.EndDate != "" {
		d = d.WhereLTE(dao.CodingTime.Columns.CreatedAt, req.EndDate)
	}
	if err = d.Scan(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}
