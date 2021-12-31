package service

// @Author: 陈健航
// @Date: 2021/2/25 20:14
// @Description:

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/gogf/gf/container/gset"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"scnu-coding/app/dao"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
)

type dockerIDEService struct {
	availablePort *gset.IntSet
}

// newDockerIDEService 构造函数
// @return s
// @date 2021-03-06 22:28:41
func newDockerIDEService() (s *dockerIDEService) {
	s = new(dockerIDEService)
	s.availablePort = gset.NewIntSet(false)
	for i := 8000; i < 8100; i++ {
		s.availablePort.Add(i)
	}
	// 拉镜像
	s.pullIDEImage()
	return s
}

func (d *dockerIDEService) pullIDEImage() {
	// 下载镜像
	IDEImageNames := g.Cfg().GetStrings("ide.image.imageNames")
	existImages, err := utils.DockerUtil.ListImages(context.Background())
	if err != nil {
		panic(err)
	}
	for _, IDEImage := range IDEImageNames {
		isExist := false
		for _, existImage := range existImages {
			if existImage.RepoTags[0] == IDEImage {
				isExist = true
				break
			}
		}
		// 不存在该镜像才pull
		if !isExist {
			if err = utils.DockerUtil.ImagePull(context.Background(), IDEImage); err != nil {
				panic(err)
			}
		}
	}
}

// OpenIDE 打开一个IDE
// @receiver d *dockerIDEService
// @param ctx context.Context
// @param Id int
// @param labID int
// @return url string
// @return err error
// @date 2021-07-22 21:46:05
func (d *dockerIDEService) OpenIDE(ctx context.Context, req *define.OpenIDEReq) (url string, err error) {
	// 上锁
	ideLock.Lock()
	defer ideLock.UnLock()
	// 查看容器是否存在
	container, err := d.isContainerExist(ctx, &req.IDEIdentifier)
	if err != nil {
		return "", err
	}
	var port uint16
	// 容器本来就存在，返回访问端口
	if container != nil && container.State == "running" {
		port = container.Ports[0].PublicPort
	} else {
		// 如果容器处于异常状态，先删掉
		if container != nil {
			if err = d.removeIDE(ctx, container); err != nil {
				return "", err
			}
			container = nil
		}
		// 新启动一个容器,获得端口
		port, err = d.StartIDE(ctx, req)
		if err != nil {
			return "", err
		}
	}
	host := g.Cfg().GetString("ide.deployment.docker.host")
	url = fmt.Sprintf("%s%s:%d", "http://", host, port)
	return url, nil
}

func (d dockerIDEService) StartIDE(ctx context.Context, req *define.OpenIDEReq) (port uint16, err error) {
	// 获得语言类型
	language, err := d.getLanguageByLabId(ctx, req.LabId)
	if err != nil {
		return 0, err
	}
	imageName := getImageName(language)
	// 端口映射
	// 随机取出一个端口
	ports := d.availablePort.Pops(1)
	// 主机端口映射
	portMapping := make(map[string]string)
	portMapping[gconv.String(ports[0])] = "8080"
	// 处理路径映射
	//mountMapping := make(map[string]string)
	mountMapping, err := d.getMountMapping(ctx, req)
	if err != nil {
		return 0, err
	}
	// 装配环境变量
	env := d.getEnv(ctx, req)
	// 容器名
	containerName := fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId)
	// 标签
	label := map[string]string{"userId": gconv.String(req.UserId), "labId": gconv.String(req.LabId)}
	// 启动容器
	err = utils.DockerUtil.RunContainer(ctx, imageName, portMapping, mountMapping, env, label, containerName)
	// 当分配的端口被占用而且被占用的端口
	for err != nil && gstr.Contains(err.Error(), "port is already allocated") {
		if d.availablePort.Size() == 0 {
			return 0, gerror.NewCode(gcode.CodeNil, "all the ports are already allocated")
		} else {
			// 重新分配一个端口
			ports = d.availablePort.Pops(1)
			portMapping[gconv.String(ports[0])] = "8080"
			err = utils.DockerUtil.RunContainer(ctx, imageName, portMapping, mountMapping, env, label, containerName)
		}
	}
	// 返回端口
	return uint16(ports[0]), nil
}

func (d dockerIDEService) getEnv(_ context.Context, req *define.OpenIDEReq) (env []string) {
	env = make([]string, 0)
	// 密码
	env = append(env, "PASSWORD=12345678")
	// 用户名
	env = append(env, fmt.Sprintf("USERID=%d", req.UserId))
	// 实验id
	env = append(env, fmt.Sprintf("LABID=%d", req.LabId))
	// 下面用于是插件和后端通信的环境变量
	ip := g.Cfg().GetString("ide.container.connect.ip")
	port := g.Cfg().GetString("ide.container.connect.port")
	openPath := g.Cfg().GetString("ide.container.connect.openPath")
	openUrl := fmt.Sprintf("http://%s:%s%s", ip, port, openPath)
	env = append(env, fmt.Sprintf("OPENURL=%s", openUrl))

	endPath := g.Cfg().GetString("ide.container.connect.endPath")
	endUrl := fmt.Sprintf("http://%s:%s%s", ip, port, endPath)
	env = append(env, fmt.Sprintf("ENDURL=%s", endUrl))
	env = append(env, "DOCKER_USER=coder")
	return env
}

func (d dockerIDEService) getMountMapping(ctx context.Context, req *define.OpenIDEReq) (mountMapping map[string]string, err error) {
	mountMapping = make(map[string]string)
	workDirContainer := "/home/coder/project"
	// 工作区使用userId和labId来标识
	workDirHost := getWorkDirHostPath(ctx, &req.IDEIdentifier)
	// 映射工作目录
	mountMapping[workDirHost] = workDirContainer
	// 映射配置目录
	language, err := d.getLanguageByLabId(ctx, req.LabId)
	if err != nil {
		return nil, err
	}
	// 配置文件目录
	//configHost := fmt.Sprintf("/data/scnu_coding/%d/.config/%d", req.UserId, language)
	configHost := getConfigPath(ctx, &req.IDEIdentifier, language)
	configContainer := "/root/.local/share/code-server"
	mountMapping[configHost] = configContainer
	return mountMapping, nil
}

func (d dockerIDEService) getLanguageByLabId(ctx context.Context, labId int) (language int, err error) {
	// 找到课程
	courseId, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.CourseId)
	if err != nil {
		return 0, err
	}
	// 找语言类型
	languageType, err := dao.Course.Ctx(ctx).WherePri(courseId).Value(dao.Course.Columns.LanguageType)
	if err != nil {
		return 0, err
	}
	return languageType.Int(), nil
}

// isContainerExist 查找是否存在这个容器
// @Description
// @receiver d
// @param ctx
// @param req
// @return container
// @return err
// @date 2021-12-21 19:51:53
func (d dockerIDEService) isContainerExist(ctx context.Context, req *define.IDEIdentifier) (container *types.Container, err error) {
	// 查找该容器是否存在
	//containerName := fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId)
	containers, err := utils.DockerUtil.ListContainer(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "label", Value: fmt.Sprintf("%s=%d", "userID", req.UserId)},
			filters.KeyValuePair{Key: "label", Value: fmt.Sprintf("%s=%d", "labID", req.LabId)}),
	})
	if err != nil {
		return nil, err
	}
	if len(containers) > 0 {
		return &containers[0], nil
	}
	return nil, nil
}

func (d *dockerIDEService) RemoveIDE(ctx context.Context, req *define.CloseIDEReq) (err error) {
	container, err := d.isContainerExist(ctx, req)
	if err != nil {
		return err
	}
	// 不存在该容器
	if container != nil {
		if err = d.removeIDE(ctx, container); err != nil {
			return err
		}
	}
	return nil
}

func (d *dockerIDEService) removeIDE(ctx context.Context, container *types.Container) (err error) {
	if err = utils.DockerUtil.RemoveContainer(ctx, container.ID); err != nil {
		return err
	}
	// 归还端口
	for _, port := range container.Ports {
		d.availablePort.Add(gconv.Int(port.PublicPort))
	}
	return nil
}

//// execStopAndRemoveIDEDocker 执行删除并移除容器
//// @receiver s
//// @params userId
//// @params language_enum
//// @params labId
//// @return err
//// @date 2021-03-06 22:19:38
//func (d *dockerIDEService) execStopAndRemoveIDEDocker(ctx context.Context, req *define.CloseIDEReq) (err error) {
//	// 停止容器
//	cmd := fmt.Sprintf("docker stop myIde-%d-%d", req.UserId, req.LabId)
//	// 这里的操作时关闭容器，有时候容器因为某些原因本来就已经关闭，这时候会报错。但目的一样就不必理会
//	_, _ = utils.DeploymentSsh.ExecCmd(cmd)
//	// 删除容器
//	cmd = fmt.Sprintf("docker rm myIde-%d-%d", req.UserId, req.LabId)
//	// 不handle error的原因同上
//	_, _ = utils.DeploymentSsh.ExecCmd(cmd)
//	return nil
//}
//
//// execRunIDEDocker 真正启动一个docker容器，注意分清楚ctx的userid和userid的区别，当学生打开自己的某个工作目录时二者一样，但当
//// 教师打开ide检查学生代码时，ide容器应该是属于教师的，ide-name也是属于教师的，但被挂载的目录是要用学生的userid确认路径
//// @receiver receiver
//// @params ctx 这里记录的才是操作人
//// @params userId 这里主要用于挂载目录
//// @params languageEnum
//// @params labId
//// @return port
//// @return err
//// @date 2021-05-08 23:40:56
//func (d *dockerIDEService) execRunIDEDocker(ctx context.Context, req *define.OpenIDEReq) (port int, err error) {
//	// 得到可用端口
//	port, err = execGetAvailablePort()
//	if err != nil {
//		return 0, err
//	}
//	languageEnum := 0
//	if req.LabId > 0 {
//		courseId, err := dao.Lab.Ctx(ctx).Cache(0).WherePri(req.LabId).Value(dao.Lab.Columns.CourseId)
//		if err != nil {
//			return 0, err
//		}
//		languageType, err := dao.Course.Ctx(ctx).Cache(0).WherePri(courseId).Value(dao.Course.Columns.LanguageType)
//		if err != nil {
//			return 0, err
//		}
//		languageEnum = languageType.Int()
//	} else {
//		languageEnum = -req.LabId
//	}
//	// 镜像地址
//	imageName := getImageName(languageEnum)
//	// 是否可编辑
//	isEditAble := ""
//	if req.IsEditAble {
//		isEditAble = "-u root"
//	}
//	// 挂载路径
//	mountedWorkspaceLocal := getWorkspacePathMounted(strconv.Itoa(req.MountedUserId), strconv.Itoa(req.LabId))
//	// 环境路径
//	mountEnvLocal := getWorkspacePathMounted(strconv.Itoa(req.UserId), fmt.Sprintf(".env-%d", languageEnum))
//	// 容器内的环境路径
//	mountEnvDocker := getDockerEnvMount(languageEnum)
//	ip := "localhost"
//	//ip, err := service.Common.GetIp()
//	if err != nil {
//		return 0, err
//	}
//	memoryLimit := g.Cfg().GetString("ide.config.memoryLimit")
//	cmd := fmt.Sprintf(
//		// 设端口
//		"docker run -itd "+
//			// 内存限制
//			"-m %s "+
//			"--init "+
//			// 端口初始化
//			"-p %d:3000 "+
//			// 用户ID
//			"-e USERID=%d "+
//			// 实验id
//			"-e LABID=%d "+
//			// 回传的地址
//			"-e BACKEND_URL=%s "+
//			// 关掉的地址
//			"-e SHUTDOWN_URL=%s "+
//			// 工作目录挂载
//			"-v %s:/home/project "+
//			// 环境目录挂载
//			"-v %s:%s "+
//			// 是否可编辑
//			"%s "+
//			// 命名，例如myIde-56-12,，56是userId,12是labId
//			"--name=myIde-%d-%d "+
//			// image
//			"%s",
//		// 内存限制
//		memoryLimit,
//		// 外置端口
//		port,
//		req.UserId,
//		req.LabId,
//		// 回访地址
//		ip+g.Cfg().GetString("server.RealAddress")+"/web/ide/alive",
//		ip+g.Cfg().GetString("server.RealAddress")+"/web/ide/alive",
//		// 挂载目录
//		mountedWorkspaceLocal,
//		// 环境目录
//		mountEnvLocal,
//		// docker里的环境目录
//		mountEnvDocker,
//		// 是否可编辑
//		isEditAble,
//		// 环境目录
//		req.UserId, req.LabId,
//		// image
//		imageName,
//	)
//	// 启动容器
//	if _, err = utils.DeploymentSsh.ExecCmd(cmd); err != nil {
//		return 0, err
//	}
//	return port, err
//}
//
//// execIsContainerAlive 检查该IDE容器是不是存活
//// @receiver receiver
//// @params languageEnum
//// @params userId
//// @params labId
//// @return isExist
//// @date 2021-04-17 00:39:59
//func (d *dockerIDEService) execIsContainerAlive(req *define.OpenIDEReq) (isExist bool) {
//	cmd := fmt.Sprintf("docker ps --filter name=myIde-%d-%d ", req.UserId, req.LabId)
//	output, err := utils.DeploymentSsh.ExecCmd(cmd)
//	if err != nil {
//		return false
//	}
//	// 是否存在
//	return gstr.ContainsI(output, fmt.Sprintf("myIde-%d-%d", req.UserId, req.LabId))
//}
