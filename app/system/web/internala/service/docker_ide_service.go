package service

// @Author: 陈健航
// @Date: 2021/2/25 20:14
// @Description:

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
)

type dockerIDEService struct{}

func (d *dockerIDEService) ListIDEContainerName(ctx context.Context) []string {
	ret := make([]string, 0)
	containers, err := utils.DockerUtil.ListContainer(ctx, map[string]string{"name": "ide"})
	if err != nil {
		return nil
	}
	for _, c := range containers {
		ret = append(ret, c.Names[0])
	}
	return ret
}

func (d *dockerIDEService) ListServerInfo(ctx context.Context) (err error) {
	serverInfo, err := utils.DockerUtil.ServerInfo(ctx)
	if err != nil {
		return err
	}
	println(serverInfo)
	return nil
}

// newDockerIDEService 构造函数
// @return s
// @date 2021-03-06 22:28:41
func newDockerIDEService() (s *dockerIDEService) {
	s = new(dockerIDEService)
	// 拉镜像
	//s.pullIDEImage()
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
			if existImage.RepoTags != nil && existImage.RepoTags[0] == IDEImage {
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
	ip := g.Cfg().GetString("docker.ip")
	url = fmt.Sprintf("%s://%s:%d", "http", ip, port)
	return url, nil
}

func (d *dockerIDEService) StartIDE(ctx context.Context, req *define.OpenIDEReq) (uint16, error) {
	// 获得语言类型
	language, err := getLanguageByLabId(ctx, req.LabId)
	if err != nil {
		return 0, err
	}
	imageName := getImageName(language)
	// 端口映射
	portMap := getPort(ctx)
	// 处理路径映射
	binds, err := getBinds(ctx, req)
	if err != nil {
		return 0, err
	}
	// 装配环境变量
	env := getEnv(ctx, req)
	// 容器名
	containerName := fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId)
	// 启动容器
	containerId, err := utils.DockerUtil.RunDocker(ctx, imageName, portMap, binds, env, containerName)
	// 当分配的端口被占用而且被占用的端口
	if err != nil {
		return 0, err
	}
	list, err := utils.DockerUtil.ListContainer(ctx, map[string]string{"id": containerId})
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, gerror.NewCode(gcode.CodeOperationFailed, "启动异常")
	}
	// 返回端口
	return list[0].Ports[0].PublicPort, nil
}

// isContainerExist 查找是否存在这个容器
// @Description
// @receiver d
// @param ctx
// @param req
// @return container
// @return err
// @date 2021-12-21 19:51:53
func (d *dockerIDEService) isContainerExist(ctx context.Context, req *define.IDEIdentifier) (container *types.Container, err error) {
	// 查找该容器是否存在
	containerName := fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId)
	containers, err := utils.DockerUtil.ListContainer(ctx, map[string]string{"name": containerName})
	if err != nil {
		return nil, err
	}
	if len(containers) > 0 {
		return &containers[0], nil
	}
	return nil, nil
}

func (d *dockerIDEService) RemoveIDE(ctx context.Context, req *define.IDEIdentifier) (err error) {
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
	return nil
}
