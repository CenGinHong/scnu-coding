package service

import (
	"context"
	"fmt"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/utils"
	"strconv"

	"github.com/gogf/gf/frame/g"
)

// @Author: 陈健航
// @Date: 2021/5/30 22:50
// @Description:

type swarm struct{}

func newSwarmService() (s *swarm) {
	s = new(swarm)
	// 代码基础存放路径
	return s
}

func (s *swarm) removeIDE(ctx context.Context, req *define.CloseIDEReq) (err error) {
	if err = s.execStopAndRemoveTheiaDocker(req); err != nil {
		return err
	}
	return nil
}

// OpenTheia 打开一个IDE
// @receiver s *swarm
// @param ctx context.Context
// @param req *define.OpenIDEReq
// @return url string
// @return err error
// @date 2021-08-29 00:24:13
func (s *swarm) OpenTheia(ctx context.Context, req *define.OpenIDEReq) (url string, err error) {
	// 上锁
	ideLock.Lock()
	defer ideLock.UnLock()
	// 查看有没有还没关闭的容器
	port, err := getIdePort(req)
	if err != nil {
		return "", err
	}
	// 端口存在
	if port != 0 {
		// 获取端口
		g.Log().Debugf("复用已经开启的IDE容器")
	} else {
		// 之前的已经关闭,重新开一个新的容器,并存入缓存
		port, err = s.execRunTheiaDocker(ctx, req)
		if err != nil {
			return "", err
		}
		// 初始化缓存信息，置入缓存
		g.Log().Debugf("开启新的IDE")
		if err = setIdePort(req, port); err != nil {
			return "", err
		}
	}
	host := g.Cfg().GetString("ide.deployment.docker.host")
	url = fmt.Sprintf("%s%s:%d", "http://:", host, port)
	return url, nil
}

func (s *swarm) execStopAndRemoveTheiaDocker(req *define.CloseIDEReq) (err error) {
	// 删除容器
	cmd := fmt.Sprintf("docker service rm myIde-%d-%d", req.UserId, req.LabId)
	// 不handle error的原因同上
	if _, err = utils.DeploymentSsh.ExecCmd(cmd); err != nil {
		return err
	}
	return nil
}

func (s *swarm) execRunTheiaDocker(ctx context.Context, req *define.OpenIDEReq) (port int, err error) {
	// 得到可用端口
	port, err = execGetAvailablePort()
	if err != nil {
		return 0, err
	}
	// 镜像地址
	imageName := getImageName(1)
	// 是否可编辑
	isEditAble := ""
	// 是否可编辑
	if req.IsEditAble {
		isEditAble = "-u root"
	}
	// 远程挂载的nfs主机地址
	mountNfsWorkspaceIp := g.Cfg().GetString("ide.storage.host")
	// nfs主机路径
	mountNfsWorkspacePath := getWorkspacePathMounted(strconv.Itoa(req.MountedUserId), strconv.Itoa(req.LabId))
	mountNfsEnvPath := getWorkspacePathMounted(strconv.Itoa(req.UserId), fmt.Sprintf(".env-%d", 1))
	mountContainerEnvPath := getDockerEnvMount(1)
	memoryLimit := g.Cfg().GetString("ide.config.memoryLimit")

	//s := "docker service create --mount 'type=volume,src=nfs-test,dst=/home/project,volume-driver=local,volume-opt=type=nfs,volume-opt=device=%s:%s,volume-opt=o=addr=192.168.1.79,vers=4,soft,timeo=180,bg,tcp,rw' -p 8888:80 --name nginx nginx:1.12"
	cmd := fmt.Sprintf("docker run -itd "+
		"%s "+
		"--mount 'type=volume,src=nfs-workspace,dst=/home/project,"+
		"volume-driver=local,volume-opt=type=nfs,volume-opt=device=%s:%s,\""+
		"volume-opt=o=addr=%s,vers=4,soft,timeo=180,bg,tcp,rw\"' "+
		"--mount 'type=volume,src=nfs-env,dst=%s,volume-driver=local,volume-opt=type=nfs,"+
		"volume-opt=device=%s:%s,\"volume-opt=o=addr=%s,vers=4,soft,timeo=180,bg,tcp,rw\"' "+
		"-p %d:3000 %s --reserve-memory %s --name myIde-%d-%d-%d %s",
		memoryLimit,
		mountNfsWorkspaceIp,
		mountNfsWorkspacePath,
		mountNfsWorkspaceIp,
		mountNfsWorkspaceIp,
		mountContainerEnvPath,
		mountNfsEnvPath,
		mountNfsWorkspaceIp,
		port,
		isEditAble,
		memoryLimit,
		// 内存限制
		1, req.UserId, req.LabId,
		// 名字里用于做标识
		imageName,
	)
	// 启动容器
	if _, err = utils.DeploymentSsh.ExecCmd(cmd); err != nil {
		return 0, err
	}
	return port, err
}
