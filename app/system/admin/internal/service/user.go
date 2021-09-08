package service

import (
	"bytes"
	"context"
	"encoding/csv"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"golang.org/x/crypto/bcrypt"
	"scnu-coding/app/dao"
	"scnu-coding/app/model"
	"scnu-coding/app/service"
	"scnu-coding/app/system/admin/internal/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/response"
	"strconv"
)

var SysUser = userService{}

type userService struct {
}

func (s *userService) GetAllUser(ctx context.Context) (resp *response.PageResp, err error) {
	// 获取分页信息
	pageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.SysUserResp, 0)
	// 筛选集
	filter := make(map[string][]*response.FilterType, 0)
	// 查找所有可筛选项 major
	majorValue, err := dao.SysUser.Ctx(ctx).Distinct().FindArray(dao.SysUser.Columns.Major)
	if err != nil {
		return nil, err
	}
	tempFilter := make([]*response.FilterType, 0)
	for _, value := range majorValue {
		tempFilter = append(tempFilter, &response.FilterType{
			Text:     value.String(),
			Value:    value.String(),
			Children: nil,
		})
	}
	filter[dao.SysUser.Columns.Major] = tempFilter
	// 查找所有可筛选项roleId
	roleValue := make([]*model.SysRole, 0)
	if err = dao.SysRole.Ctx(ctx).Distinct().Fields(dao.SysRole.Columns.RoleId, dao.SysRole.Columns.Description).Scan(&roleValue); err != nil {
		return nil, err
	}
	tempFilter = make([]*response.FilterType, 0)
	for _, role := range roleValue {
		tempFilter = append(tempFilter, &response.FilterType{
			Text:     role.Description,
			Value:    strconv.Itoa(role.RoleId),
			Children: nil,
		})
	}
	filter["role"] = tempFilter
	// 查找真正数据项
	d := dao.SysUser.Ctx(ctx).Page(pageInfo.Current, pageInfo.PageSize)
	// 筛选项
	for k, v := range pageInfo.ParseFilterFields {
		d = d.Where(k, v)
	}
	// 查总数
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	// 排序
	d = d.Order(pageInfo.SortField, pageInfo.SortOrder)
	// 查询
	if err = d.Fields(define.SysUserResp{}).WithAll().Scan(&records); err != nil {
		return nil, err
	}

	resp = response.GetPageResp(records, total, filter)
	return resp, nil
}

func (s *userService) UpdateUser(ctx context.Context, req define.UpdateSysUserReq) (err error) {
	if count, err := dao.SysUser.Ctx(ctx).WhereNot(dao.SysUser.Columns.UserId, req.UserId).
		Where(dao.SysUser.Columns.Email, req.Email).Count(); err != nil {
		return err
	} else if count > 0 {
		return gerror.New("邮箱已经被使用，请更换邮箱")
	}
	if count, err := dao.SysUser.Ctx(ctx).WhereNot(dao.SysUser.Columns.UserId, req.UserId).
		Where(dao.SysUser.Columns.UserNum, req.UserNum).Count(); err != nil {
		return err
	} else if count > 0 {
		return gerror.New("学号/职工号已被其他用户使用，请检查")
	}
	if _, err = dao.SysUser.Ctx(ctx).WherePri(req.UserId).Data(req).Update(); err != nil {
		return err
	}
	return nil
}

func (s *userService) ResetPassword(ctx context.Context, Id int) (err error) {
	// 找出这个人的邮箱
	email, err := dao.SysUser.Ctx(ctx).WherePri(Id).Value(dao.SysUser.Columns.Email)
	if err != nil {
		return err
	}
	// 将邮箱作为他的新密码
	hashedPassword, err := bcrypt.GenerateFromPassword(email.Bytes(), bcrypt.DefaultCost)
	if _, err = dao.SysUser.Ctx(ctx).WherePri(Id).Data(dao.SysUser.Columns.Password, hashedPassword).Update(); err != nil {
		return err
	}
	return nil
}

func (s userService) GetStudentImportTemplate(_ context.Context) (templateFile *bytes.Buffer, err error) {
	templateFile = &bytes.Buffer{}
	//写入bom头
	utils.WriteBom(templateFile)
	w := csv.NewWriter(templateFile)
	defer w.Flush()
	header := make([]string, 0)
	header = append(header, "姓名*")
	header = append(header, "学号*")
	header = append(header, "性别")
	header = append(header, "邮箱（若不填默认使用学院邮箱）")
	header = append(header, "单位*")
	header = append(header, "专业*")
	header = append(header, "密码")
	if err = w.Write(header); err != nil {
		return nil, err
	}
	return templateFile, nil
}

func (s userService) ImportStudent(ctx context.Context, template *ghttp.UploadFile, isOverWrite bool) (err error) {
	file, err := template.Open()
	if err != nil {
		return err
	}
	file1, err := utils.RemoveBom(file)
	if err != nil {
		return err
	}
	reader := csv.NewReader(file1)
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}
	userDatas := make([]*define.ImportStudent, 0)
	for i, row := range rows {
		if i == 0 {
			continue
		}
		userNum := row[1]
		// 是否覆写已存在得用户，如果不覆写查到数据库有这条就直接跳过
		if !isOverWrite {
			count, err := dao.SysUser.Ctx(ctx).WherePri(dao.SysUser.Columns.UserNum, userNum).Count()
			if err != nil {
				return err
			}
			if count > 0 {
				continue
			}
		}
		var gender int
		if row[2] == "男" {
			gender = 1
		} else if row[2] == "女" {
			gender = 2
		}
		password := row[6]
		var hashedPassword []byte
		if password == "" {
			//默认以学号做密码
			password = userNum
		}
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		userData := &define.ImportStudent{
			Username:     row[0],
			UserNum:      userNum,
			Gender:       gender,
			Email:        row[3],
			Organization: row[4],
			Major:        row[5],
			Password:     string(hashedPassword),
		}
		userDatas = append(userDatas, userData)
	}
	if _, err = dao.SysUser.Ctx(ctx).Data(userDatas).Batch(len(userDatas)).Insert(); err != nil {
		return err
	}
	return nil
}
