package service

// @Author: 陈健航
// @Date: 2021/2/25 20:14
// @Description:

import (
	"context"
	"fmt"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/utils"
	"strconv"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
)

type dockerTheiaService struct{}

// NewDockerTheiaService 构造函数
// @return s
// @date 2021-03-06 22:28:41
func newDockerTheiaService() (s *dockerTheiaService) {
	s = new(dockerTheiaService)
	// 代码基础存放路径
	return s
}

// OpenTheia 打开一个IDE
// @receiver d *dockerTheiaService
// @param ctx context.Context
// @param Id int
// @param labID int
// @return url string
// @return err error
// @date 2021-07-22 21:46:05
func (d *dockerTheiaService) OpenTheia(ctx context.Context, req *define.OpenIDEReq) (url string, err error) {
	// 上锁
	ideLock.Lock()
	defer ideLock.UnLock()
	// 查看有没有还没关闭的容器
	port, err := getIdePort(req)
	if err != nil {
		return "", err
	}
	// 缓存存在，还要进行二次检查，docker容器不存在自愈，防止容器崩溃后处在exit状态无法启动新的
	if port != 0 && d.execIsContainerAlive(req) {
		// 获取端口
		g.Log().Debugf("复用已经开启的IDE容器")
	} else {
		// 关闭之前可能开启的容器
		if err = d.execStopAndRemoveTheiaDocker(ctx, &define.CloseIDEReq{
			IDEIdentifier: define.IDEIdentifier{
				UserId:       req.UserId,
				LanguageEnum: req.LanguageEnum,
				LabId:        req.LabId,
			},
		}); err != nil {
			return "", err
		}
		// 之前的已经关闭,重新开一个新的容器,并存入缓存
		port, err = d.execRunTheiaDocker(req)
		if err != nil {
			return "", err
		}
		g.Log().Debugf("开启新的IDE")
		// 初始化缓存信息，置入缓存
		if err = setIdePort(req, port); err != nil {
			return "", err
		}
	}
	host := g.Cfg().GetString("ide.deployment.docker.host")
	url = fmt.Sprintf("%s%s:%d", "http://", host, port)
	return url, nil
}

func (d *dockerTheiaService) removeIDE(ctx context.Context, req *define.CloseIDEReq) (err error) {
	if err = d.execStopAndRemoveTheiaDocker(ctx, req); err != nil {
		return err
	}
	return nil
}

// execStopAndRemoveTheiaDocker 执行删除并移除容器
// @receiver s
// @params userId
// @params language
// @params labId
// @return err
// @date 2021-03-06 22:19:38
func (d *dockerTheiaService) execStopAndRemoveTheiaDocker(ctx context.Context, req *define.CloseIDEReq) (err error) {
	// 停止容器
	cmd := fmt.Sprintf("docker stop myIde-%d-%d-%d", req.LanguageEnum, req.UserId, req.LabId)
	// 这里的操作时关闭容器，有时候容器因为某些原因本来就已经关闭，这时候会报错。但目的一样就不必理会
	_, _ = utils.DeploymentSsh.ExecCmd(cmd)
	// 删除容器
	cmd = fmt.Sprintf("docker rm myIde-%d-%d-%d", req.LanguageEnum, req.UserId, req.LabId)
	// 不handle error的原因同上
	_, _ = utils.DeploymentSsh.ExecCmd(cmd)
	return nil
}

// execRunTheiaDocker 真正启动一个docker容器，注意分清楚ctx的userid和userid的区别，当学生打开自己的某个工作目录时二者一样，但当
// 教师打开ide检查学生代码时，ide容器应该是属于教师的，ide-name也是属于教师的，但被挂载的目录是要用学生的userid确认路径
// @receiver receiver
// @params ctx 这里记录的才是操作人
// @params userId 这里主要用于挂载目录
// @params languageEnum
// @params labId
// @return port
// @return err
// @date 2021-05-08 23:40:56
func (d *dockerTheiaService) execRunTheiaDocker(req *define.OpenIDEReq) (port int, err error) {
	// 得到可用端口
	port, err = execGetAvailablePort()
	if err != nil {
		return 0, err
	}
	// 镜像地址
	imageName := getImageName(req.LanguageEnum)
	// 是否可编辑
	isEditAble := ""
	if req.IsEditAble {
		isEditAble = "-u root"
	}
	// 挂载路径
	mountWorkSpace := getWorkspacePathRemote(strconv.Itoa(req.MountedUserId), strconv.Itoa(req.LabId))
	// 环境路径
	mountEnvLocal := getWorkspacePathRemote(strconv.Itoa(req.MountedUserId), fmt.Sprintf(".env-%d", req.LanguageEnum))
	// 容器内的环境路径
	mountEnvDocker := getDockerEnvMount(req.LanguageEnum)
	memoryLimit := g.Cfg().GetString("ide.image.memoryLimit")
	cmd := fmt.Sprintf(
		// 设端口
		"docker run -itd --memory %s --init -p %d:3000 "+
			"%s "+
			"-v %s:/home/project "+
			// 工作目录挂载
			"-v %s:%s "+
			// 环境目录挂载
			"--name=myIde-%d-%d-%d "+
			// 命名，例如myIde-2-56-12,，56是userId,12是labId
			"%s",
		// 语言版本
		memoryLimit,
		// 内存限制
		port,
		isEditAble,
		mountWorkSpace,
		// 主机上的工作目录
		mountEnvLocal,
		// ide的环境目录
		mountEnvDocker,
		// docker里的环境目录
		req.LanguageEnum, req.UserId, req.LabId,
		// 名字里用于做标识
		imageName,
		// ide的名称
	)
	// 启动容器
	if _, err = utils.DeploymentSsh.ExecCmd(cmd); err != nil {
		return 0, err
	}
	return port, err
}

// execIsContainerAlive 检查该IDE容器是不是存活
// @receiver receiver
// @params languageEnum
// @params userId
// @params labId
// @return isExist
// @date 2021-04-17 00:39:59
func (d *dockerTheiaService) execIsContainerAlive(req *define.OpenIDEReq) (isExist bool) {
	cmd := fmt.Sprintf("docker ps --filter name=myIde-%d-%d-%d ", req.LanguageEnum, req.UserId, req.LabId)
	output, err := utils.DeploymentSsh.ExecCmd(cmd)
	if err != nil {
		return false
	}
	// 是否存在
	return gstr.ContainsI(output, fmt.Sprintf("myIde-%d-%d-%d", req.LanguageEnum, req.UserId, req.LabId))
}
