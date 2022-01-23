package service

// @Author: 陈健航
// @Date: 2021/2/25 20:14
// @Description:

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/go-connections/nat"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
)

type dockerIDEService struct{}

// newDockerIDEService 构造函数
// @return s
// @date 2021-03-06 22:28:41
func newDockerIDEService() (s *dockerIDEService) {
	s = new(dockerIDEService)
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
	host := g.Cfg().GetString("ide.deployment.docker.host")
	url = fmt.Sprintf("%s%s:%d", "http://", host, port)
	return url, nil
}

func (d dockerIDEService) StartIDE(ctx context.Context, req *define.OpenIDEReq) (uint16, error) {
	// 获得语言类型
	language, err := getLanguageByLabId(ctx, req.LabId)
	if err != nil {
		return 0, err
	}
	imageName := getImageName(language)
	// 端口映射
	portMap, err := d.getPort(ctx)
	if err != nil {
		return 0, nil
	}
	// 处理路径映射
	binds, err := d.getBinds(ctx, req)
	if err != nil {
		return 0, err
	}
	// 装配环境变量
	env := d.getEnv(ctx, req)
	// 容器名
	containerName := fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId)
	// 启动容器
	container, err := utils.DockerUtil.RunContainer(ctx, imageName, portMap, binds, env, containerName)
	// 当分配的端口被占用而且被占用的端口
	if err != nil {
		return 0, err
	}
	list, err := utils.DockerUtil.ListContainer(ctx, types.ContainerListOptions{Filters: filters.NewArgs(filters.KeyValuePair{Key: "id", Value: container.ID})})
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, gerror.NewCode(gcode.CodeOperationFailed, "启动异常")
	}
	// 返回端口
	return list[0].Ports[0].PublicPort, nil
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
	env = append(env, fmt.Sprintf("CONNECT_URL=%s", openUrl))
	env = append(env, "DOCKER_USER=coder")
	return env
}

func (d dockerIDEService) getBinds(ctx context.Context, req *define.OpenIDEReq) (mountMapping []string, err error) {
	mountMapping = make([]string, 0)
	// 工作区使用userId和labId来标识
	workDirHost := getWorkDirHostPath(ctx, &req.IDEIdentifier)
	// 映射工作目录
	mountMapping = append(mountMapping, fmt.Sprintf("%s:/home/coder/project", workDirHost))
	// 映射配置目录
	configHost, err := getConfigPath(ctx, &req.IDEIdentifier)
	if err != nil {
		return nil, err
	}
	mountMapping = append(mountMapping, fmt.Sprintf("%s:/root/.local/share/code-server", configHost))
	return mountMapping, nil
}

func (d dockerIDEService) getPort(_ context.Context) (nat.PortMap, error) {
	// 端口映射
	portMap := make(nat.PortMap, 0)
	port, err := nat.NewPort("tcp", "8080")
	if err != nil {
		return nil, err
	}
	portBind := nat.PortBinding{}
	tmp := make([]nat.PortBinding, 0, 1)
	tmp = append(tmp, portBind)
	portMap[port] = tmp
	return portMap, nil
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
	containerName := fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId)
	containers, err := utils.DockerUtil.ListContainer(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: containerName}),
	})
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
