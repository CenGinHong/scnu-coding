package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/swarm"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"scnu-coding/app/system/web/internala/define"
	"scnu-coding/app/utils"
)

// @Author: 陈健航
// @Date: 2021/5/30 22:50
// @Description:

type swarmIDEService struct{}

func (s *swarmIDEService) ListIDEContainerName(ctx context.Context) []string {
	ret := make([]string, 0)
	svces, err := utils.DockerUtil.ListService(ctx, map[string]string{"name": "ide"})
	if err != nil {
		return nil
	}
	for _, svce := range svces {
		ret = append(ret, svce.Spec.Name)
	}
	return ret
}

func newSwarmService() (s *swarmIDEService) {
	s = new(swarmIDEService)
	// 代码基础存放路径
	return s
}

func (s *swarmIDEService) RemoveIDE(ctx context.Context, req *define.IDEIdentifier) (err error) {
	// 根据名字找id
	listService, err := utils.DockerUtil.ListService(ctx, map[string]string{
		"name": fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId),
	})
	if err != nil {
		return err
	}
	if len(listService) == 0 {
		return nil
	}
	// 用id删除
	if err = utils.DockerUtil.RemoveService(ctx, listService[0].ID); err != nil {
		return err
	}
	return nil
}

func (s *swarmIDEService) ListServerInfo(ctx context.Context) (err error) {
	//TODO implement me
	panic("implement me")
}

// OpenIDE
// @Description
// @receiver s
// @param ctx
// @param req
// @return url
// @return err
// @date 2022-03-22 10:15:08
func (s *swarmIDEService) OpenIDE(ctx context.Context, req *define.OpenIDEReq) (url string, err error) {
	// 上锁
	ideLock.Lock()
	defer ideLock.UnLock()
	// 查看容器是否存在
	service, err := s.isServiceExist(ctx, &req.IDEIdentifier)
	if err != nil {
		return "", err
	}
	var port uint32
	// 容器本来就存在，返回访问端口
	if service != nil {
		port = service.Endpoint.Ports[0].PublishedPort
	} else {
		// 新启动一个容器,获得端口
		port, err = s.StartIDE(ctx, req)
		if err != nil {
			return "", err
		}
	}
	ip := g.Cfg().GetString("docker.ip")
	url = fmt.Sprintf("%s://%s:%d", "http", ip, port)
	return url, nil
}

func (s *swarmIDEService) isServiceExist(ctx context.Context, req *define.IDEIdentifier) (svce *swarm.Service, err error) {
	serviceName := fmt.Sprintf("ide-%d-%d", req.UserId, req.LabId)
	services, err := utils.DockerUtil.ListService(ctx, map[string]string{"name": serviceName})
	if err != nil {
		return nil, err
	}
	if len(services) > 0 {
		return &services[0], nil
	}
	return nil, nil
}

func (s *swarmIDEService) StartIDE(ctx context.Context, req *define.OpenIDEReq) (port uint32, err error) {
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
	containerId, err := utils.DockerUtil.RunService(ctx, imageName, portMap, binds, env, containerName)
	if err != nil {
		return 0, err
	}
	list, err := utils.DockerUtil.ListService(ctx, map[string]string{"id": containerId})
	if err != nil {
		return 0, err
	}
	if len(list) == 0 {
		return 0, gerror.NewCode(gcode.CodeOperationFailed, "启动异常")
	}
	// 返回端口
	return list[0].Endpoint.Ports[0].PublishedPort, nil
}
