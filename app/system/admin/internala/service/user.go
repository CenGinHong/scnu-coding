package service

import (
	"bufio"
	"bytes"
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gvalid"
	"github.com/xuri/excelize/v2"
	"golang.org/x/crypto/bcrypt"
	"mime/multipart"
	"scnu-coding/app/dao"
	"scnu-coding/app/model"
	"scnu-coding/app/service"
	"scnu-coding/app/system/admin/internala/define"
	"scnu-coding/library/response"
	"strconv"
	"strings"
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

func (s *userService) ListUser(ctx context.Context) (resp *response.PageResp, err error) {
	// 获取分页信息
	pageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*model.SysUser, 0)
	d := dao.SysUser.Ctx(ctx)

	// 筛选集
	filter := make(map[string][]*response.FilterType, 0)
	// 查找所有可筛选项 gender, grade, school, organization,major
	filterFields := []string{dao.SysUser.Columns.Gender,
		dao.SysUser.Columns.School, dao.SysUser.Columns.Organization, dao.SysUser.Columns.Major}
	for _, fields := range filterFields {
		array, err := d.Distinct().FindArray(fields)
		if err != nil {
			return nil, err
		}
		tempFilter := make([]*response.FilterType, 0)
		for _, value := range array {
			tempFilter = append(tempFilter, &response.FilterType{
				Text:     value.String(),
				Value:    value.String(),
				Children: nil,
			})
		}
		filter[fields] = tempFilter
	}

	// 排序
	if pageInfo != nil {
		if pageInfo.SortOrder != "" {
			d = d.Order(pageInfo.SortField, pageInfo.SortOrder)
		}
		// 加入筛选项
		for key, value := range pageInfo.ParseFilterFields {
			d = d.Where(key, value)
		}
		d = d.Page(pageInfo.Current, pageInfo.PageSize)
	}
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if err = d.Scan(&records); err != nil {
		return nil, err
	}

	resp = response.GetPageResp(records, total, filter)
	return resp, nil
}

func (s *userService) GetUser(ctx context.Context, userId int) (resp *model.SysUser, err error) {
	resp = &model.SysUser{}
	if err = dao.SysUser.Ctx(ctx).WherePri(userId).Scan(&resp); err != nil {
		return nil, err
	}
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

//func (s userService) GetStudentImportTemplate(_ context.Context) (templateFile *bytes.Buffer, err error) {
//	templateFile = &bytes.Buffer{}
//	//写入bom头
//	utils.WriteBom(templateFile)
//	w := csv.NewWriter(templateFile)
//	defer w.Flush()
//	header := make([]string, 0)
//	header = append(header, "姓名*")
//	header = append(header, "学号*")
//	header = append(header, "性别")
//	header = append(header, "邮箱（若不填默认使用学院邮箱）")
//	header = append(header, "单位*")
//	header = append(header, "专业*")
//	header = append(header, "密码")
//	if err = w.Write(header); err != nil {
//		return nil, err
//	}
//	return templateFile, nil
//}

func (s *userService) ImportStudent(ctx context.Context, template *ghttp.UploadFile, roleId int) (errMsg string, err error) {
	file, err := template.Open()
	if err != nil {
		return "", err
	}
	defer func(file multipart.File) {
		if err = file.Close(); err != nil {
			glog.Error(err)
		}
	}(file)
	reader, err := excelize.OpenReader(bufio.NewReader(file))
	if err != nil {
		return "", err
	}
	rows, err := reader.GetRows("Sheet1")
	if err != nil {
		return "", err
	}
	rows = rows[1:]
	insertData := make([]*define.ImportStudent, 0)
	sb := &strings.Builder{}
	for _, r := range rows {
		username := r[0]
		userNum := r[1]
		// 是否覆写已存在得用户，如果不覆写查到数据库有这条就直接跳过
		userId, err := dao.SysUser.Ctx(ctx).Where(dao.SysUser.Columns.UserNum, userNum).Value(dao.SysUser.Columns.UserId)
		if err != nil {
			return "", err
		}
		if !userId.IsNil() {
			continue
		}
		email := r[2]
		school := r[3]
		major := r[4]
		organization := r[5]
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userNum), bcrypt.DefaultCost)
		if err != nil {
			return "", err
		}
		userData := &define.ImportStudent{
			RoleId:       roleId,
			Username:     username,
			UserNum:      userNum,
			Email:        email,
			Organization: organization,
			Major:        major,
			School:       school,
			Password:     string(hashedPassword),
		}
		if err = gvalid.CheckStruct(ctx, userData, nil); err != nil {
			sb.WriteString(err.Error())
			sb.WriteString("/n")
			continue
		}
		insertData = append(insertData, userData)
	}
	if len(insertData) > 0 {
		if _, err = dao.SysUser.Ctx(ctx).Data(insertData).Batch(len(insertData)).Insert(); err != nil {
			return "", err
		}
	}
	return sb.String(), nil
}

func (s *userService) GetImportDemoCsv(_ context.Context) (buffer *bytes.Buffer, err error) {
	f := excelize.NewFile()
	_ = f.NewSheet("Sheet1")
	defer func(f *excelize.File) {
		if err := f.Close(); err != nil {
			glog.Error(err)
		}
	}(f)
	header := []string{"姓名", "学号", "邮箱", "学院", "专业", "班级"}
	if err = f.SetSheetRow("Sheet1", "A1", &header); err != nil {
		return nil, err
	}
	buffer, err = f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (s *userService) DeleteUser(ctx context.Context, userId int) (err error) {
	if _, err = dao.SysUser.Ctx(ctx).WherePri(userId).Delete(); err != nil {
		return err
	}
	return nil
}
