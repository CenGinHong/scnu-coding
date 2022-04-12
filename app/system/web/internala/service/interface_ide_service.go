package service

// @Author: 陈健航
// @Date: 2021/3/3 23:31
// @Description:

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"scnu-coding/app/dao"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
	"scnu-coding/library/enum/language_enum"
)

type iIDE interface {
	// OpenIDE 启动或复用一个IDE容器
	OpenIDE(ctx context.Context, req *define.OpenIDEReq) (url string, err error)
	// RemoveIDE 关闭容器操作
	RemoveIDE(ctx context.Context, req *define.IDEIdentifier) (err error)
	// ListServerInfo 列举服务器信息
	ListServerInfo(ctx context.Context) (err error)

	ListIDEContainerName(ctx context.Context) []string
}

//// idePortCache 记录每一个容器所占用
//var idePortCache = utils.NewMyCache()

// 锁
var ideLock = utils.NewMyMutex()

// newIDE 新建IDE管理服务
// @return t
// @date 2021-05-03 00:04:22
func newIDE() (t iIDE) {
	// 获取配置文件中填入的配置方式
	switch g.Cfg().GetString("ide.deploymentType") {
	case "docker":
		t = newDockerIDEService()
	case "swarm":
		t = newSwarmService()
	default:
		panic("不支持的IDE容器部署方式")
	}
	return t
}

// getImageName 获取容器名
// @params languageEnum
// @return imageName
// @date 2021-05-03 00:04:41
func getImageName(languageEnum int) (imageName string) {
	imageNames := g.Cfg().GetStrings("ide.image.imageNames")
	languageString := ""
	switch languageEnum {
	case language_enum.Full:
		languageString = language_enum.Num2LanguageString(language_enum.Full)
	case language_enum.Cpp:
		languageString = language_enum.Num2LanguageString(language_enum.Cpp)
	case language_enum.Java:
		languageString = language_enum.Num2LanguageString(language_enum.Java)
	case language_enum.Python:
		languageString = language_enum.Num2LanguageString(language_enum.Python)
	}
	// 找出对应的镜像返回
	for _, imageName = range imageNames {
		if gstr.ContainsI(imageName, languageString) {
			return imageName
		}
	}
	return ""
}

func getIDEWorkDirHostPath(_ context.Context, ident *define.IDEIdentifier) (workDirPath string) {
	workspaceBasePathRemote := g.Cfg().GetString("ide.storage.workspaceBasePathRemote")
	workDirPath = fmt.Sprintf("%s/%d/%d", workspaceBasePathRemote, ident.UserId, ident.LabId)
	return workDirPath
}

func getServiceLocalPath(_ context.Context, ident *define.IDEIdentifier) (workDirPath string) {
	localPath := g.Cfg().GetString("ide.storage.serviceLocalPath")
	workDirPath = fmt.Sprintf("%s/%d/%d", localPath, ident.UserId, ident.LabId)
	return workDirPath
}

func getIDEConfigPath(ctx context.Context, ident *define.IDEIdentifier) (configPath string, err error) {
	configBasePathRemote := g.Cfg().GetString("ide.storage.configBasePathRemote")
	language, err := getLanguageByLabId(ctx, ident.LabId)
	if err != nil {
		return "", err
	}
	configPath = fmt.Sprintf("%s/%d/.config/%d", configBasePathRemote, ident.UserId, language)
	return configPath, err
}

func getLanguageByLabId(ctx context.Context, labId int) (language int, err error) {
	// 找到课程
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return 0, err
	}
	// 找语言类型
	languageType, err := dao.Course.Ctx(ctx).WherePri(dao.Course.Columns.CourseId, courseId).Value(dao.Course.Columns.LanguageType)
	if err != nil {
		return 0, err
	}
	return languageType.Int(), nil
}

func getEnv(_ context.Context, req *define.OpenIDEReq) (env []string) {
	env = make([]string, 0)
	// 密码
	env = append(env, "PASSWORD=12345678")
	// 用户名
	env = append(env, fmt.Sprintf("USERID=%d", req.UserId))
	// 实验id
	env = append(env, fmt.Sprintf("LABID=%d", req.LabId))
	// 下面用于是插件和后端通信的环境变量
	ip := g.Cfg().GetString("ide.container.heartbeat.ip")
	port := g.Cfg().GetString("ide.container.heartbeat.port")
	heartbeatPath := g.Cfg().GetString("ide.container.heartbeat.heartbeatPath")
	openUrl := fmt.Sprintf("http://%s:%s%s", ip, port, heartbeatPath)
	env = append(env, fmt.Sprintf("CONNECT_URL=%s", openUrl))
	env = append(env, "DOCKER_USER=coder")
	return env
}

func getBinds(ctx context.Context, req *define.OpenIDEReq) (mountMapping []string, err error) {
	mountMapping = make([]string, 0)
	// 工作区使用userId和labId来标识
	workDirHost := getIDEWorkDirHostPath(ctx, &req.IDEIdentifier)
	// 映射工作目录
	workDirMapping := fmt.Sprintf("%s:/home/coder/project", workDirHost)
	if !req.IsEditAble {
		workDirMapping = workDirMapping + ":ro"
	}
	mountMapping = append(mountMapping, workDirMapping)
	// 映射配置目录
	configHost, err := getIDEConfigPath(ctx, &req.IDEIdentifier)
	if err != nil {
		return nil, err
	}
	mountMapping = append(mountMapping, fmt.Sprintf("%s:/root/.local/share/code-server", configHost))
	return mountMapping, nil
}

func getPort(_ context.Context) (portMap []string) {
	// 端口映射
	portMap = []string{"0:8080"}
	return portMap
}
