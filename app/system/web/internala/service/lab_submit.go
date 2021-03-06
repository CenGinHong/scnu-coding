package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"github.com/xuri/excelize/v2"
	"scnu-coding/app/dao"
	"scnu-coding/app/model"
	"scnu-coding/app/service"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/enum/language_enum"
	"scnu-coding/library/response"
	"strings"
	"time"
)

// @Author: 陈健航
// @Date: 2021/3/8 16:55
// @Description:

// LabSummit service
var LabSummit = labSummitService{}

type labSummitService struct{}

// UpdateReport 保存实验报告
// @receiver l *labSummitService
// @param ctx context.Context
// @param req *define.UpdateReportContentReq
// @return err error
// @date 2021-07-29 19:44:25
func (l *labSummitService) UpdateReport(ctx context.Context, req *define.UpdateReportContentReq) (err error) {
	ctxUser := service.Context.Get(ctx).User
	if _, err = dao.LabSubmit.Ctx(ctx).Data(g.Map{
		dao.LabSubmit.Columns.UserId:        ctxUser.UserId,
		dao.LabSubmit.Columns.LabId:         req.LabId,
		dao.LabSubmit.Columns.ReportContent: req.ReportContent,
	}).Save(); err != nil {
		return err
	}
	return nil
}

// GetReportContent 查找实验报告
// @receiver l *labSummitService
// @param ctx context.Context
// @param labId int
// @return reportContent string
// @return err error
// @date 2021-07-29 19:43:46
func (l *labSummitService) GetReportContent(ctx context.Context, req *define.GetReportContentReq) (reportContent string, err error) {
	ctxUser := service.Context.Get(ctx).User
	if req.UserId == 0 {
		req.UserId = ctxUser.UserId
	}
	value, err := dao.LabSubmit.Ctx(ctx).Where(g.Map{
		dao.LabSubmit.Columns.UserId: req.UserId,
		dao.LabSubmit.Columns.LabId:  req.LabId,
	}).Value(dao.LabSubmit.Columns.ReportContent)
	if err != nil {
		return "", err
	}
	if value.IsNil() {
		reportContent = ""
	} else {
		reportContent = value.String()
	}
	return reportContent, nil
}

// GetSubmitCode 教师检查代码用
// @receiver l *labSummitService
// @param ctx context.Context
// @param labSubmitId int
// @return resp *define.GetReportContentAndCodeResp
// @return err error
// @date 2021-08-12 23:18:26
func (l *labSummitService) GetSubmitCode(ctx context.Context, labSubmitId int) (resp *define.GetReportContentAndCodeResp, err error) {
	labSubmit := &model.LabSubmit{}
	value, err := dao.LabSubmit.Ctx(ctx).WherePri(labSubmitId).Fields(dao.LabSubmit.Columns.ReportContent, dao.LabSubmit.Columns.UserId, dao.LabSubmit.Columns.LabId).One()
	if err != nil {
		return nil, err
	}
	if err = value.Struct(&labSubmit); err != nil {
		return nil, err
	}
	resp = &define.GetReportContentAndCodeResp{}
	resp.ReportContent = labSubmit.ReportContent
	if resp.Code, err = l.GetCodeData(ctx, &define.ReadCodeDataReq{
		LabId:  labSubmit.LabId,
		UserId: labSubmit.UserId,
	}); err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateFinishStat 是否已经完成编写代码
// @receiver s
// @params req
// @return err
// @date 2021-03-08 17:08:06
func (l *labSummitService) UpdateFinishStat(ctx context.Context, req *define.UpdateLabFinishReq) (err error) {
	ctxUser := service.Context.Get(ctx).User
	if req.UserId == 0 {
		req.UserId = ctxUser.UserId
	}
	if _, err = dao.LabSubmit.Ctx(ctx).Data(g.Map{
		dao.LabSubmit.Columns.LabId:    req.LabId,
		dao.LabSubmit.Columns.UserId:   req.UserId,
		dao.LabSubmit.Columns.IsFinish: req.IsFinish,
	}).Save(); err != nil {
		return err
	}
	return nil
}

func (l *labSummitService) ListLabSummit(ctx context.Context, labId int) (resp *response.PageResp, err error) {
	ctxPageInfo := service.Context.Get(ctx).PageInfo
	records := make([]*define.ListLabSubmitResp, 0)
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return nil, err
	}
	// 找出所有选课的学生
	d := dao.ReCourseUser.Ctx(ctx)
	d = d.Where(dao.ReCourseUser.Columns.CourseId, courseId)
	total, err := d.Count()
	if err != nil {
		return nil, err
	}
	if ctxPageInfo != nil {
		d = d.Page(ctxPageInfo.Current, ctxPageInfo.PageSize)
	}
	if err = d.WithAll().Scan(&records); err != nil {
		return nil, err
	}
	if err = dao.LabSubmit.Ctx(ctx).
		Where(g.Map{
			dao.LabSubmit.Columns.UserId: gdb.ListItemValuesUnique(records, "UserId"),
			dao.LabSubmit.Columns.LabId:  labId,
		}).
		Fields(define.ListLabSubmitResp{}.LabSubmitDetail).
		ScanList(&records, "LabSubmitDetail", "user_id:UserId"); err != nil {
		return nil, err
	}
	for _, r := range records {
		r.IsIDERunning = IDE.IsIDERunning(r.UserId, labId)
	}
	resp = response.GetPageResp(records, total, nil)
	return resp, nil
}

func (l *labSummitService) ListLabSummitId(ctx context.Context, labId int) (resp []*define.ListLabSubmitIdResp, err error) {
	resp = make([]*define.ListLabSubmitIdResp, 0)
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return nil, err
	}
	// 找出所有选课的学生
	d := dao.ReCourseUser.Ctx(ctx)
	d = d.Where(dao.ReCourseUser.Columns.CourseId, courseId)
	if err != nil {
		return nil, err
	}
	if err = d.WithAll().Scan(&resp); err != nil {
		return nil, err
	}
	if err = dao.LabSubmit.Ctx(ctx).Where(dao.LabSubmit.Columns.LabId, labId).Where(dao.LabSubmit.Columns.UserId,
		gdb.ListItemValuesUnique(resp, "UserId")).Fields(define.ListLabSubmitResp{}.LabSubmitDetail).
		ScanList(&resp, "LabSubmitDetail", "user_id:UserId"); err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateScoreAndComment 实验打分评语
// @receiver receiver
// @params req
// @return err
// @date 2021-04-21 10:19:20
func (l *labSummitService) UpdateScoreAndComment(ctx context.Context, req *define.UpdateLabSummitScoreAndCommentReq) (err error) {
	if _, err = dao.LabSubmit.Ctx(ctx).Data(g.Map{
		dao.LabSubmit.Columns.LabSubmitComment: req.Comment,
		dao.LabSubmit.Columns.Score:            req.Score,
		dao.LabSubmit.Columns.UserId:           req.UserId,
		dao.LabSubmit.Columns.LabId:            req.LabId,
	}).Save(); err != nil {
		return err
	}
	return nil
}

// ExecPlagiarismCheckByMoss 执行moss
// @params basePath
// @params languageEnum
// @params userIds
// @params labId
// @return url
// @return err
// @date 2021-04-17 00:47:40
func (l *labSummitService) ExecPlagiarismCheckByMoss(ctx context.Context, labId int) (url string, err error) {
	// 查出课程
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return "", err
	}
	// 查出语言类型
	languageType, err := dao.Course.Ctx(ctx).WherePri(courseId).Value(dao.Course.Columns.LanguageType)
	if err != nil {
		return "", err
	}
	var language string
	var ext []string
	switch languageType.Int() {
	case language_enum.Cpp:
		language = "cc"
		ext = append(ext, "*.cpp", "*.h", "*.c")
	case language_enum.Java:
		language = "java"
		ext = append(ext, "*.java")
	case language_enum.Python:
		language = "python"
		ext = append(ext, "*.py")
	default:
		return "", gerror.NewCode(gcode.CodeNotSupported, "暂不支持的语言类型")
	}
	// 查出选课学生
	userIds, err := dao.ReCourseUser.Ctx(ctx).WherePri(dao.ReCourseUser.Columns.CourseId, courseId).Array(dao.ReCourseUser.Columns.UserId)
	if err != nil {
		return "", err
	}
	// 新建所有的
	mossClient, err := utils.NewMossClient(language, g.Cfg().GetString("moss.userId"))
	if err != nil {
		return "", err
	}
	defer func(mossClient *utils.MossClient) {
		_ = mossClient.Close()
	}(mossClient)
	if err = mossClient.Run(); err != nil {
		return "", err
	}
	// 本地放置代码的位置
	uploadFilePaths := make([]string, 0)
	for _, userId := range userIds {
		// 该学生的实验工作目录，注意该目录可能未创建(因为学生还没有开始做实验）
		path := getServiceLocalPath(ctx, &define.IDEIdentifier{
			UserId: gconv.Int(userId),
			LabId:  labId,
		})
		// 不存在该目录，不理了
		if !gfile.Exists(path) {
			continue
		}
		// 读出所有文件路径
		filePath, err := gfile.ScanDirFile(path, gstr.Join(ext, ","), true)
		if err != nil {
			return "", err
		}
		// 加入比对
		uploadFilePaths = append(uploadFilePaths, filePath...)
	}
	// 上传所有的文件
	for _, uploadFilePath := range uploadFilePaths {
		if err = mossClient.UploadFile(uploadFilePath, false); err != nil {
			return "", err
		}
	}
	// 关闭
	if err = mossClient.SendQuery(); err != nil {
		return "", err
	}
	res := mossClient.ResultURL
	return res.String(), err
}

//// CollectCompilerErrorLog 收集编译错误报告
//// @receiver receiver
//// @params labId
//// @return compilerErrorLog
//// @return err
//// @date 2021-04-21 10:18:56
//func (l *labSummitService) CollectCompilerErrorLog(labId int) (compilerErrorLog string, err error) {
//	// 找出该门课的学生
//	courseId, err := dao.Lab.WherePri(labId).Cache(0).Value(dao.Lab.Columns.CourseId)
//	if err != nil {
//		return "", err
//	}
//	// 所有选了这门课的学生id
//	stuIds, err := g.Model(dao.ReCourseUser.Table).Where(dao.ReCourseUser.Columns.CourseId, courseId).Array(dao.ReCourseUser.Columns.UserId)
//	if err != nil {
//		return "", err
//	}
//	// 这个是后端服务中挂载的代码存储主机上的路径
//	// 新建channel
//	compilerErrorLogChan := make(chan string, len(stuIds))
//	for _, stuId := range stuIds {
//		logFileName := getWorkspacePathLocal(stuId.String(), strconv.Itoa(labId), ".compilerErrorLog")
//		// 开goroutine收集所有协程
//		go func() {
//			content := gfile.GetContentsWithCache(logFileName)
//			compilerErrorLogChan <- content
//		}()
//	}
//	// 收集结果
//	sb := &strings.Builder{}
//	for i := 0; i < len(stuIds); i++ {
//		log := <-compilerErrorLogChan
//		if log != "" {
//			sb.WriteString(log)
//		}
//	}
//	// 返回结果
//	return sb.String(), nil
//}
//
//// PlagiarismCheck 代码查重检测
//// @receiver receiver
//// @params labId
//// @return resp
//// @return err
//// @date 2021-04-21 10:18:44
//func (receiver *labSummitService) PlagiarismCheck(labId int) (resp []*model.PlagiarismCheckResp, err error) {
//	// 找出课程id
//	courseId, err := dao.Lab.WherePri(labId).Cache(0).FindValue(dao.Lab.Columns.CourseId)
//	if err != nil {
//		return nil, nil
//	}
//	// 找出userId
//	stuIds, err := g.Table(dao.ReCourseUser.Table).Where(dao.ReCourseUser.Columns.CourseId, courseId).FindArray(dao.ReCourseUser.Columns.UserID)
//	if err != nil {
//		return nil, err
//	}
//	// 找出语言类型
//	languageEnum, err := dao.Course.WherePri(courseId).FindValue(dao.Course.Columns.Language)
//	if err != nil {
//		return nil, err
//	}
//
////// 执行moss查重
////url, err := receiver.ExecPlagiarismCheckByMoss(languageEnum.Int(), stuIds, labId)
////if err != nil {
////	return nil, nil
////}
////// 解析结果
////resp, err = receiver.parsePlagiarismCheck(url)
////if err != nil {
////	return nil, err
////}
////// 装填字段
////userDetail := make([]*model.SysUser, 0)
////if err = dao.SysUser.WherePri(stuIds).
////	Fields(dao.SysUser.Columns.RealName, dao.SysUser.Columns.UserID, dao.SysUser.Columns.Num).
////	FindScan(&userDetail); err != nil {
////	return nil, nil
////}
////for _, v := range resp {
////	finish := 0
////	for _, v1 := range userDetail {
////		if v.UserId1 == v1.UserID {
////			v.RealName1 = v1.RealName
////			v.Num1 = v1.Num
////			finish++
////		}
////		if v.UserId2 == v1.UserID {
////			v.RealName2 = v1.RealName
////			v.Num2 = v1.Num
////			finish++
////		}
////		if finish == 2 {
////			break
////		}
////	}
////}
////	return resp, err
////}

// GetCodeData 生成代码树，用户快速查看代码
// @receiver receiver
// @params req
// @return resp
// @return err
// @date 2021-04-17 00:42:16
func (l *labSummitService) GetCodeData(ctx context.Context, req *define.ReadCodeDataReq) (resp []*define.CodeData, err error) {
	// 只查出这几种类型的代码文件
	extNames := []string{"*.txt", "*.py", "*.java", "*.cpp", "*.c", "*.h", "*.ts", "*.js", "Makefile", ".json"}
	pathPrefix := getServiceLocalPath(ctx, &define.IDEIdentifier{UserId: gconv.Int(req.UserId), LabId: gconv.Int(req.LabId)})
	glog.Debug(pathPrefix)
	type TempCodeFile struct {
		Filename string // 文件名
		Content  string // 文件内容
	}
	retChan := make(chan *TempCodeFile)
	root := make([]*define.CodeData, 0)
	// 路径未创建，说明学生未打开ide做实验
	if !gfile.Exists(pathPrefix) {
		return root, nil
	}
	// 读取该同学该实验工作目录的全部文件
	scanFiles, err := gfile.ScanDirFileFunc(pathPrefix, strings.Join(extNames, ","), true, func(path string) string {
		// 用goroutine读取文件
		go func() {
			tempCodeFile := &TempCodeFile{
				Filename: path,
				Content:  gfile.GetContentsWithCache(path, 5*time.Second),
			}
			// 用channel返回
			retChan <- tempCodeFile
		}()
		return path
	})
	if err != nil {
		return nil, err
	}
	// 读取遍历出来的文件
	for i := 0; i < len(scanFiles); i++ {
		tempCodeFile := <-retChan
		fileName := tempCodeFile.Filename
		// 去掉前缀
		fileName = gstr.TrimStr(fileName, pathPrefix+"/")
		// 切割目录级
		fileNameSplit := gstr.Split(fileName, "/")
		glog.Debug(fileNameSplit)
		// 构造树
		l.buildTreeNode(&root, fileNameSplit, 0, tempCodeFile.Content)
	}
	return root, nil
}

// buildTreeNode 构建代码树结构
// @Description
// @receiver l
// @param childNode 这里在切片扩容需要使用是指针传递
// @param path
// @param index
// @param content
// @date 2022-03-02 10:18:12
func (l *labSummitService) buildTreeNode(childNode *[]*define.CodeData, path []string, index int, content string) {
	glog.Debug(path)
	// 到达叶子节点，该节点是一个文件
	if index == len(path)-1 {
		*childNode = append(*childNode, &define.CodeData{
			Name:      path[index],
			Content:   content,
			ChildNode: make([]*define.CodeData, 0),
		})
		return
	}
	// 是否存在该层目录
	isExist := false
	for _, child := range *childNode {
		if child.Name == path[index] {
			// 递归构建树
			l.buildTreeNode(&child.ChildNode, path, index+1, content)
			isExist = true
			break
		}
	}
	// 不存在，创建该层目录
	if !isExist {
		newNode := &define.CodeData{
			ChildNode: make([]*define.CodeData, 0),
			Name:      path[index],
		}
		// 挂载叶子节点
		*childNode = append(*childNode, newNode)
		// 递归构建树
		l.buildTreeNode(&newNode.ChildNode, path, index+1, content)
	}
}

//func (receiver *labSummitService) parsePlagiarismCheck(url string) (records []*model.PlagiarismCheckResp, err error) {
//	htmlResp, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer func(Body io.ReadCloser) {
//		_ = Body.Close()
//	}(htmlResp.Body)
//	root, err := htmlquery.Parse(htmlResp.Body)
//	if err != nil {
//		return nil, err
//	}
//	tr := htmlquery.Find(root, "/html/body/table/tbody/tr[*]")
//	records = make([]*model.PlagiarismCheckResp, 0)
//	resChannel := make(chan *model.PlagiarismCheckResp, len(tr)-1)
//	for i, v := range tr {
//		if i == 0 {
//			continue
//		}
//		// 开多个协程爬取信息
//		go func(node *html.Node, index int, resChannel chan *model.PlagiarismCheckResp) {
//			ret := &model.PlagiarismCheckResp{}
//			// 通过channel返回
//			defer func() {
//				resChannel <- ret
//			}()
//			// 解析html
//			detailUrl := node.FirstChild.FirstChild.Attr[0].Val
//			dir1 := node.FirstChild.FirstChild.FirstChild.Data
//			dir2 := node.FirstChild.NextSibling.FirstChild.FirstChild.Data
//			workspaceBasePathLocal := g.Cfg().GetString("ide.storage.workspaceBasePathLocal")
//			length := gstr.PosI(gstr.TrimStr(dir1, path.Join(workspaceBasePathLocal, "codespaces")), "/", 1) - 1
//			userId1 := gstr.SubStr(gstr.TrimStr(dir1, path.Join(workspaceBasePathLocal, "codespaces")), 1, length)
//			length = gstr.PosI(gstr.TrimStr(dir2, path.Join(workspaceBasePathLocal, "codespaces")), "/", 1) - 1
//			userId2 := gstr.SubStr(gstr.TrimStr(dir2, path.Join(workspaceBasePathLocal, "codespaces")), 1, length)
//			length = gstr.PosI(dir1, "%") - gstr.PosI(dir1, "(") - 1
//			similarity := gstr.SubStr(dir1, gstr.PosI(dir1, "(")+1, length)
//			ret.Url = detailUrl
//			ret.UserId1 = gconv.Int(userId1)
//			ret.UserId2 = gconv.Int(userId2)
//			ret.Similarity = gconv.Int(similarity)
//		}(v, i, resChannel)
//	}
//
//	// 收集返回信息
//	for i := 0; i < len(tr)-1; i++ {
//		select {
//		case ret := <-resChannel:
//			records = append(records, ret)
//		}
//	}
//	return records, nil
//}

func (l *labSummitService) ExportScore(ctx context.Context, req *define.ExportLabScoreReq) (buffer *bytes.Buffer, err error) {
	// 学生名单
	stuRecords := make([]*define.ExportLabScore, 0)
	// 找出所有学生名单
	if err = dao.ReCourseUser.Ctx(ctx).
		Where(dao.ReCourseUser.Columns.CourseId, req.CourseId).
		WithAll().
		Scan(&stuRecords); err != nil {
		return nil, err
	}
	// 绑定成绩详情
	if err = dao.LabSubmit.Ctx(ctx).Where(
		g.Map{
			dao.LabSubmit.Columns.LabId:  req.LabIds,
			dao.LabSubmit.Columns.UserId: gdb.ListItemValues(stuRecords, "UserId"),
		}).
		Fields(define.ExportLabScore{}.LabSubmitDetails).
		OrderAsc(dao.LabSubmit.Columns.LabId).
		ScanList(&stuRecords, "LabSubmitDetails", "user_id:UserId"); err != nil {
		return nil, err
	}
	f := excelize.NewFile()
	_ = f.NewSheet("成绩")
	defer func(f *excelize.File) {
		if err = f.Close(); err != nil {
			glog.Error(err)
		}
	}(f)
	labRecords := make([]model.Lab, 0)
	if err = dao.Lab.Ctx(ctx).WherePri(req.LabIds).
		Fields(dao.Lab.Columns.LabId, dao.Lab.Columns.Title).
		OrderAsc(dao.Lab.Columns.LabId).
		Scan(&labRecords); err != nil {
		return nil, err
	}
	header := []string{"姓名", "学号"}
	for _, r := range labRecords {
		header = append(header, r.Title)
	}
	if err = f.SetSheetRow("成绩", "A1", &header); err != nil {
		return nil, err
	}
	for i, stuRecord := range stuRecords {
		row := make([]interface{}, 0)
		row = append(row, stuRecord.UserDetail.Username)
		row = append(row, stuRecord.UserDetail.UserNum)
		stuRecordIdx := 0
		if stuRecord.LabSubmitDetails == nil {
			// 防止数组为空
			stuRecord.LabSubmitDetails = make([]*struct {
				LabId  int `orm:"lab_id"`
				Score  int `orm:"score"`
				UserId int `orm:"user_id"`
			}, 0)
		}
		for _, record := range labRecords {
			if stuRecordIdx < len(stuRecord.LabSubmitDetails) && stuRecord.LabSubmitDetails[stuRecordIdx].LabId == record.LabId {
				row = append(row, stuRecord.LabSubmitDetails[stuRecordIdx].Score)
				stuRecordIdx++
			} else {
				row = append(row, "未评分")
			}
		}
		// 从A2开始填起
		if err = f.SetSheetRow("成绩", fmt.Sprintf("A%d", i+2), &row); err != nil {
			return nil, err
		}
	}
	if buffer, err = f.WriteToBuffer(); err != nil {
		return nil, err
	}
	return buffer, nil
}
