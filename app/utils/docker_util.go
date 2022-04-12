package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/util/gconv"
	"io"
	"scnu-coding/app/system/web/internala/define"
	"strings"
	"time"
)

var DockerUtil = newDockerUtil()

type dockerUtil struct {
	client *client.Client
}

func newDockerUtil() (d dockerUtil) {
	// 拼装docker remote api 地址
	host := fmt.Sprintf("tcp://%s:%s", g.Cfg().GetString("docker.ip"),
		g.Cfg().GetString("docker.port"))
	isTslVerify := g.Cfg().GetBool("docker.withTlsVerify")
	if isTslVerify {
		// 连接加密
		path := g.Cfg().GetString("docker.ca")
		cacertPath := fmt.Sprintf("%s/%s", path, "ca.pem")
		certPath := fmt.Sprintf("%s/%s", path, "cert.pem")
		keyPath := fmt.Sprintf("%s/%s", path, "key.pem")
		cli, err := client.NewClientWithOpts(client.WithHost(host),
			client.WithTLSClientConfig(cacertPath, certPath, keyPath))
		if err != nil {
			panic(err)
		}
		d = dockerUtil{cli}
	} else {
		cli, err := client.NewClientWithOpts(client.WithHost(host), client.WithAPIVersionNegotiation())
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
func (d *dockerUtil) ListImages(ctx context.Context) (imageList []types.ImageSummary, err error) {
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
func (d *dockerUtil) GetContainerStat(ctx context.Context, containerID string) (containerStat *define.ContainerStat, err error) {
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
func (d *dockerUtil) ListContainer(ctx context.Context, filter map[string]string) (containers []types.Container, err error) {
	args := make([]filters.KeyValuePair, 0)
	for k, v := range filter {
		args = append(args, filters.KeyValuePair{Key: k, Value: v})
	}
	containers, err = d.client.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filters.NewArgs(args...),
	})
	if err != nil {
		return nil, err
	}
	return containers, nil
}

// RunDocker
// @Description
// @receiver d
// @param ctx
// @param imageName 镜像名
// @param portMap 端口映射，形如 ["80:8080"]
// @param binds 路径映射，形如["/home/test:/home"]
// @param env 环境变量
// @param containerName 镜像名
// @return *container.ContainerCreateCreatedBody
// @return error
// @date 2022-03-20 17:08:58
func (d *dockerUtil) RunDocker(ctx context.Context, imageName string,
	portMap []string, binds []string, env []string, containerName string) (containerId string, err error) {
	// 准备端口
	p := make(nat.PortMap, 0)
	for _, s := range portMap {
		// 拆分为宿主机端口和容器端口
		split := strings.Split(s, ":")
		// 绑定容器端口
		port, err := nat.NewPort("tcp", split[1])
		if err != nil {
			return "0", err
		}
		portBind := nat.PortBinding{}
		if split[0] != "0" {
			portBind.HostPort = split[0]
		}
		tmp := make([]nat.PortBinding, 0, 1)
		tmp = append(tmp, portBind)
		p[port] = tmp
	}
	if err != nil {
		return "", err
	}
	// 构建容器
	c, err := d.client.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Env:   env,
		Tty:   true,
		User:  "root",
	}, &container.HostConfig{
		Binds:        binds,
		PortBindings: p,
		Privileged:   true,
		AutoRemove:   true,
	}, nil, nil, containerName)
	if err != nil {
		return "", err
	}
	// 启动容器
	if err = d.client.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return c.ID, nil
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

func (d *dockerUtil) DaemonHost() string {
	return d.client.DaemonHost()
}

func (d *dockerUtil) RunService(ctx context.Context, imageName string,
	portMap []string, binds []string, env []string, containerName string) (serviceId string, err error) {
	// 准备nfs挂载
	nfsAddr := g.Cfg().GetString("ide.storage.nfsAddr")
	mountList := make([]mount.Mount, 0)
	for _, bind := range binds {
		split := strings.Split(bind, ":")
		mountList = append(mountList, mount.Mount{Type: mount.TypeVolume, Target: split[1], VolumeOptions: &mount.VolumeOptions{
			DriverConfig: &mount.Driver{
				Name: "local",
				Options: map[string]string{
					"type":   "nfs",
					"o":      fmt.Sprintf("addr=%s,rw", nfsAddr),
					"device": fmt.Sprintf(":%s", split[0]),
				},
			},
		}})
	}
	// 准备端口
	p := make([]swarm.PortConfig, 0)
	for _, s := range portMap {
		// 拆分为宿主机端口和容器端口
		split := strings.Split(s, ":")
		// 绑定容器端口
		tmp := swarm.PortConfig{TargetPort: gconv.Uint32(split[1])}
		if split[0] != "0" {
			tmp.PublishedPort = gconv.Uint32(split[0])
		}
		p = append(p, tmp)
	}
	// 创建服务
	serviceCreateResponse, err := d.client.ServiceCreate(ctx, swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: containerName,
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image:  imageName,
				User:   "root",
				Env:    env,
				TTY:    true,
				Mounts: mountList,
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: p,
		},
	}, types.ServiceCreateOptions{})
	if err != nil {
		return "0", err
	}
	return serviceCreateResponse.ID, nil
}

func (d *dockerUtil) ListService(ctx context.Context, filter map[string]string) (services []swarm.Service, err error) {
	args := make([]filters.KeyValuePair, 0)
	for k, v := range filter {
		args = append(args, filters.KeyValuePair{Key: k, Value: v})
	}
	services, err = d.client.ServiceList(ctx, types.ServiceListOptions{Filters: filters.NewArgs(args...)})
	if err != nil {
		return nil, err
	}
	return services, nil
}

func (d *dockerUtil) RemoveService(ctx context.Context, id string) (err error) {
	if err = d.client.ServiceRemove(ctx, id); err != nil {
		return err
	}
	return nil
}

func (d *dockerUtil) T() {
	_, err := d.client.ServiceCreate(context.Background(), swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: "demo",
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: &swarm.ContainerSpec{
				Image: "nginx",
				User:  "root",
				TTY:   true,
				Mounts: []mount.Mount{{Type: mount.TypeVolume, Target: "/home", VolumeOptions: &mount.VolumeOptions{
					DriverConfig: &mount.Driver{
						Name: "local",
						Options: map[string]string{
							"type":   "nfs",
							"o":      fmt.Sprintf("addr=%s,rw", "10.50.3.213"),
							"device": fmt.Sprintf(":%s", "/home/horace/testNfs4"),
						},
					},
				}}},
			},
		},
		EndpointSpec: &swarm.EndpointSpec{
			Ports: []swarm.PortConfig{{TargetPort: 80}},
		},
	}, types.ServiceCreateOptions{})
	if err != nil {
		return
	}

	list, err := d.client.ServiceList(context.Background(), types.ServiceListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{Key: "name", Value: "demo"}),
	})
	//info, err := d.ServerInfo(context.Background())
	//if err != nil {
	//	return
	//}
	//nodes, err := d.client.NodeList(context.Background(), types.NodeListOptions{})
	//if err != nil {
	//	return
	//}
	//
	//if err != nil {
	//	return
	//}
	d.client.ServiceRemove(context.Background(), list[0].ID)
}
