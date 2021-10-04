package service

// @Author: 陈健航
// @Date: 2021/3/3 23:31
// @Description:

import (
	"context"
	"fmt"
	"math/rand"
	"path"
	"scnu-coding/app/system/web/internal/define"
	"scnu-coding/app/utils"
	"sync"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

type iTheia interface {
	// OpenTheia 启动或复用一个IDE容器
	OpenTheia(ctx context.Context, req *define.OpenIDEReq) (url string, err error)
	// removeIDE 关闭容器操作
	removeIDE(ctx context.Context, req *define.CloseIDEReq) (err error)
}

// idePortCache 记录每一个容器所占用
var idePortCache = utils.NewMyCache()

// 锁
var ideLock = utils.NewMyMutex()

// newIdeService 新建IDE管理服务
// @return t
// @date 2021-05-03 00:04:22
func newTheia() (t iTheia) {
	// 获取配置文件中填入的配置方式
	switch g.Cfg().GetString("ide.deploymentType") {
	case "docker":
		t = newDockerTheiaService()
	//case "k3s":
	//	t = newK3sTheiaService()
	case "swarm":
		//t = newSwarmService()
	default:
		panic("不支持的IDE容器部署方式")
	}
	// 在服务重启阶段清理所有已开启的IDE容器
	//clearAllIde(t.removeIDE)
	return t
}

// getImageName 获取容器名
// @params languageEnum
// @return imageName
// @date 2021-05-03 00:04:41
func getImageName(languageEnum int) (imageName string) {
	switch languageEnum {
	case 0:
		imageName = g.Cfg().GetString("ide.image.name.full")
	case 1:
		imageName = g.Cfg().GetString("ide.image.name.cpp")
	case 2:
		imageName = g.Cfg().GetString("ide.image.name.java")
	case 3:
		imageName = g.Cfg().GetString("ide.image.name.python")
	}
	return imageName
}

// getDockerEnvMount IDE容器的环境目录（每个容器都不太一样）
// @params languageEnum
// @return environmentMount
// @date 2021-05-03 00:04:53
func getDockerEnvMount(languageEnum int) (environmentMount string) {
	switch languageEnum {
	case 0:
		//full
		environmentMount = "/home/ide/.ide"
	case 1:
		//cpp
		environmentMount = "/root/.ide"
	case 2:
		//java
		environmentMount = "/root/.ide"
	case 3:
		//python
		environmentMount = "/home/ide/.ide"
	case 4:
		//web
		environmentMount = "/root/.ide"
	}
	return environmentMount
}

// removeIdePort 删除IDE标识块
// @params languageEnum
// @params userId
// @params labId
// @return err
// @date 2021-04-17 00:46:24
func removeIdePort(userId int, labId int) (err error) {
	key := fmt.Sprintf("%d-%d", userId, labId)
	// 删除stat
	if _, err = idePortCache.Remove(key); err != nil {
		return err
	}
	return nil
}

// getIdePort 获取IDE标识块
// @params languageEnum
// @params userId
// @params labId
// @return stat
// @return err
// @date 2021-05-03 00:05:31
func getIdePort(req *define.OpenIDEReq) (port int, err error) {
	key := fmt.Sprintf("%d-%d", req.UserId, req.LabId)
	data, err := idePortCache.GetVar(key)
	if err != nil {
		return 0, err
	}
	if data.IsNil() {
		return 0, nil
	}
	// 有未关闭的容器，返回端口
	return data.Int(), nil
}

// setIdePort 设置缓存标识块
// @params languageEnum
// @params userId
// @params labId
// @params value
// @return err
// @date 2021-04-17 00:47:21
func setIdePort(req *define.OpenIDEReq, port int) (err error) {
	key := fmt.Sprintf("%d-%d", req.UserId, req.LabId)
	if err = idePortCache.Set(key, port, 0); err != nil {
		return err
	}
	return nil
}

// execGetAvailablePort 获取可用端口
// @return randPort
// @return err
// @date 2021-05-03 00:05:44
func execGetAvailablePort() (randPort int, err error) {
	rand.Seed(time.Now().UnixNano())
	for {
		// 生成30000-32000的一个随机数
		randPort = rand.Intn(2000) + 30000
		cmd := fmt.Sprintf("netstat -an | grep %d", randPort)
		// 这里不处理error了，有输出会伴随错误,没有输出代表该接口可用
		if output, _ := utils.DeploymentSsh.ExecCmd(cmd); output == "" {
			break
		}
	}
	return randPort, nil
}

func clearAllIde(removeIdeF func(userId int, labId int) (err error)) {
	// 上锁
	ideLock.Lock()
	defer ideLock.UnLock()
	output, err := utils.DeploymentSsh.ExecCmd("docker ps -a --format \"table {{.Names}} \"")
	if err != nil {
		g.Log().Error(err)
	}
	containerNames := gstr.Split(output, "\n")[1:]
	wg := &sync.WaitGroup{}
	for _, containerName := range containerNames {
		if gstr.HasPrefix(containerName, "myIde") {
			wg.Add(1)
			go func(key string) {
				defer wg.Done()
				gstr.TrimStr(key, "myIde-")
				split := gstr.Split(key, "-")
				Id, labId := gconv.Int(split[0]), gconv.Int(split[1])
				if err = removeIdeF(Id, labId); err != nil {
					g.Log().Error(err)
				}
				// 移除缓存
				_ = removeIdePort(Id, labId)
			}(containerName)
		}
	}
	wg.Wait()
}

func getWorkspacePathLocal(params ...string) (workspacePath string) {
	workspaceBasePathLocal := g.Cfg().GetString("ide.storage.workspaceBasePathLocal")
	workspacePath = path.Join(workspaceBasePathLocal, "codespaces")
	for _, param := range params {
		workspacePath = path.Join(workspacePath, param)
	}
	return workspacePath
}

func getWorkspacePathMounted(params ...string) (workspacePath string) {
	workspaceBasePathRemote := g.Cfg().GetString("ide.storage.workspaceBasePathRemote")
	workspacePath = path.Join(workspaceBasePathRemote, "codespaces")
	for _, param := range params {
		workspacePath = path.Join(workspacePath, param)
	}
	return workspacePath
}

// moss不用脚本了
//// initEnv 初始化环境
//// @date 2021-04-17 00:41:58
//func initEnv() {
//	workspaceBasePathLocal := g.Cfg().GetString("ide.workspaceBasePathLocal")
//	// 把moss的脚本下载下来
//	mossUrl := g.Cfg().GetString("moss.downLoad")
//	mossAccount := g.Cfg().GetString("moss.userId")
//	// 创建目录
//	cmd := fmt.Sprintf("mkdir -p %s/util && "+
//		// 下载moss脚本
//		"cd %s/util && if [ ! -e moss ]; then wget %s && "+
//		// 替换userid
//		"sed -i '167c $userid=%s;' moss && chmod a+x %s/util; fi",
//		// 赋予执行权限
//		workspaceBasePathLocal, workspaceBasePathLocal, mossUrl, mossAccount, workspaceBasePathLocal)
//	if _, err := gproc.ShellExec(cmd); err != nil {
//		panic("initEnv" + err.Error())
//	}
//}
