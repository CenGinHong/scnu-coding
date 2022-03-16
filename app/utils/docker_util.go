package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"io"
	"scnu-coding/app/system/web/internala/define"
	"time"
)

var DockerUtil = newDockerUtil()

type dockerUtil struct {
	client *client.Client
}

func newDockerUtil() (d dockerUtil) {
	// 拼装docker remote api 地址
	host := g.Cfg().GetString("ide.deployment.docker.host")
	port := g.Cfg().GetString("ide.deployment.docker.port")
	protocol := g.Cfg().GetString("ide.deployment.docker.protocol")
	addr := fmt.Sprintf("%s://%s:%s", protocol, host, port)
	if protocol == "tcp" {
		cli, err := client.NewClientWithOpts(client.WithHost(addr), client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}
		d = dockerUtil{cli}
	} else {
		// tls加密连接
		caPath := g.Cfg().GetString("ide.ide.deployment.docker.tls.caPath")
		certPath := g.Cfg().GetString("ide.ide.deployment.docker.tls.certPath")
		keyPath := g.Cfg().GetString("ide.ide.deployment.docker.tls.keyPath")
		cli, err := client.NewClientWithOpts(client.WithHost(addr), client.WithAPIVersionNegotiation(),
			client.WithTLSClientConfig(caPath, certPath, keyPath))
		if err != nil {
			panic(err)
		}
		d = dockerUtil{cli}
	}
	return d
}

// ListImages
// @Description 列出所有的容器镜像
// @receiver d
// @param ctx
// @return imageList
// @return err
// @date 2021-12-21 10:56:50
func (d dockerUtil) ListImages(ctx context.Context) (imageList []types.ImageSummary, err error) {
	imageList, err = d.client.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}
	return imageList, nil
}

// GetContainerStat 获得容器的运行时状态
// @Description
// @receiver d
// @param ctx
// @param containerID
// @return containerStat
// @return err
// @date 2021-12-30 13:18:00
func (d dockerUtil) GetContainerStat(ctx context.Context, containerID string) (containerStat *define.ContainerStat, err error) {
	stats, err := d.client.ContainerStatsOneShot(ctx, containerID)
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(stats.Body)
	buf := make([]byte, 1024)
	statsByte := make([]byte, 0)
	for {
		n, err := stats.Body.Read(buf)
		statsByte = append(statsByte, buf[:n]...)
		//表示读取完毕
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
	}
	containerStat = &define.ContainerStat{}
	if err = json.Unmarshal(statsByte, containerStat); err != nil {
		return nil, err
	}
	return containerStat, nil
}

// ListContainer
// @Description 列出所有容器
// @receiver d
// @param ctx
// @param opts
// @return containers
// @return err
// @date 2021-12-21 10:57:48
func (d *dockerUtil) ListContainer(ctx context.Context, opts types.ContainerListOptions) (containers []types.Container, err error) {
	containers, err = d.client.ContainerList(ctx, opts)

	if err != nil {
		return nil, err
	}
	return containers, nil
}

// RunContainer
// @Description 运行一个容器
// @receiver d
// @param ctx
// @param imageName 容器名称
// @param portMap 目录挂载匹配，前面的是宿主机卷，后面是容器端口
// @param binds 端口匹配，全面是宿主机端口，后面是容器端口
// @param env 环境目录
// @param containerName
// @return error
// @date 2021-12-21 10:58:08
func (d *dockerUtil) RunContainer(ctx context.Context, imageName string,
	portMap nat.PortMap, binds []string, env []string, containerName string) (*container.ContainerCreateCreatedBody, error) {
	// 构建容器
	c, err := d.client.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
		Tty:   false,
		User:  "root",
	}, &container.HostConfig{
		Binds:        binds,
		PortBindings: portMap,
		Privileged:   true,
		AutoRemove:   true,
	}, nil, nil, containerName)
	if err != nil {
		return nil, err
	}
	// 启动容器
	if err = d.client.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	return &c, nil
}

// ImagePull
// @Description 录取容器目录
// @receiver d
// @param ctx
// @param imageName
// @return err
// @date 2021-12-21 11:05:16
func (d *dockerUtil) ImagePull(ctx context.Context, imageName string) (err error) {
	glog.Infof("Pulling image:%s......", imageName)
	_, err = d.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	glog.Infof("Pulling image:%s finish", imageName)
	return nil
}

func (d *dockerUtil) StopContainer(ctx context.Context, containerId string) (err error) {
	timeout := time.Second * 10
	if err = d.client.ContainerStop(ctx, containerId, &timeout); err != nil {
		return err
	}
	return nil
}

func (d *dockerUtil) RemoveContainer(ctx context.Context, containerId string) (err error) {
	if err = d.client.ContainerRemove(ctx, containerId, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}); err != nil {
		return err
	}
	return nil
}

func (d *dockerUtil) RestartContainer(ctx context.Context, containerId string) (err error) {
	t := 3 * time.Second
	if err = d.client.ContainerRestart(ctx, containerId, &t); err != nil {
		return err
	}
	return nil
}

func (d *dockerUtil) ServerInfo(ctx context.Context) (info *types.Info, err error) {
	info1, err := d.client.Info(ctx)
	if err != nil {
		return nil, err
	}
	return &info1, nil
}

func (d *dockerUtil) CreateVolumeForNfs(ctx context.Context, volumeName string, addr string, path string) {
	// 查找该存储卷是否已经被创建
	list, err := d.client.VolumeList(ctx, filters.NewArgs(filters.KeyValuePair{Key: "name", Value: volumeName}))
	if err != nil {
		return
	}
	// 已存在，返回
	if len(list.Volumes) > 0 {
		return
	}
	// 创建nfs存储卷
	if _, err = d.client.VolumeCreate(ctx, volume.VolumeCreateBody{
		Driver: "local",
		DriverOpts: map[string]string{
			"type":   "nfs",
			"o":      fmt.Sprintf("addr=%s,rw", addr),
			"device": fmt.Sprintf(":%s", path),
		},
	}); err != nil {
		return
	}
}

func (d dockerUtil) T() {

}
