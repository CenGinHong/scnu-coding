package service

import (
	"context"
	"scnu-coding/app/system/web/internala/define"
)

// @Author: 陈健航
// @Date: 2021/3/27 0:51
// @Description:

type k3sIDEService struct{}

func (k *k3sIDEService) OpenIDE(ctx context.Context, req *define.OpenIDEReq) (url string, err error) {
	return "", err
}

func (k *k3sIDEService) RemoveIDE(ctx context.Context, req *define.IDEIdentifier) (err error) {
	panic("implement me")
}

func newK3sIDEService() (k *k3sIDEService) {
	k = new(k3sIDEService)
	// 创建命名空间
	//if err := k.(); err != nil {
	//	panic(err)
	//}
	return k
}

//func (receiver *k3sIDEService) createNameSpace() (err error) {
//	cmd := "kubectl create ns code_platform"
//	if _, err = utils.DeploymentSsh.ExecCmd(cmd); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (receiver *k3sIDEService) removeIdE(languageEnum int, Id int, labId int) (err error) {
//	if err = receiver.execUninstallTheiaDeployment(languageEnum, Id, labId); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (receiver *k3sIDEService) execMkDir(dir string) (err error) {
//	// 因为打算用nfs把文件挂载到后端服务容器中，所以直接用了
//	cmd := fmt.Sprintf("if [ ! -d %s ]; then mkdir -p %s; fi;", dir, dir)
//	if _, err = gproc.ShellExec(cmd); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (receiver *k3sIDEService) execRunTheiaDeployment(ctx context.Context, Id int, languageEnum int, labId int) (port int, err error) {
//	ctxUser := service.Context.Get(ctx).User
//	// 得到可用端口
//	port, err = execGetAvailablePort()
//	if err != nil {
//		return
//	}
//	chartName := g.Cfg().GetString("ide.image.chartName")
//	// 镜像地址
//	imageName := getImageName(languageEnum)
//	split := gstr.Split(imageName, ":")
//	repository := split[0]
//	tag := split[1]
//	mountWorkSpaceRemote := getWorkspacePathMounted(strconv.Itoa(Id), strconv.Itoa(labId))
//	mountEnvRemote := getWorkspacePathMounted(strconv.Itoa(Id), fmt.Sprintf(".env-%d", languageEnum))
//	// 预建文件夹
//	if err = receiver.execMkDir(mountWorkSpaceRemote); err != nil {
//		return 0, err
//	}
//	if err = receiver.execMkDir(mountEnvRemote); err != nil {
//		return 0, err
//	}
//	mountEnvDocker := getDockerEnvMount(languageEnum)
//	runAsUser := 0
//	// 目前暂定的规则是，仅当用户本身为学生而且打算打开一个过期的实验时会没有写权限
//	if ctxUser.RoleId == role.STUDENT {
//		// 查询截止日期
//		ddl, err := dao.Lab.Ctx(ctx).WherePri(labId).Value(dao.Lab.Columns.Deadline)
//		if err != nil {
//			return 0, err
//		}
//		// 过了截止时间,将不可编辑
//		if (!ddl.IsNil() || !ddl.IsEmpty()) && gtime.Now().After(ddl.GTime()) {
//			runAsUser = 1000
//		}
//	}
//	cmd := fmt.Sprintf(
//		"helm install --namespace=code_platform myIde-%d-%d-%d %s "+
//			"--set service.nodePort=%d "+
//			// 端口
//			"--set image.repository=%s "+
//			// 仓库名
//			"--set image.tag=%s "+
//			//标签
//			"--set image.volumeMounts.env=%s "+
//			// 容器内环境目录挂载
//			"--set volumes.nfs.workspace.server=%s "+
//			// 代码储存主机
//			"--set volumes.nfs.workspace.path=%s "+
//			// 代码主机工作目录挂载地址
//			"--set volumes.nfs.env.server=%s "+
//			// 代码储存主机
//			"--set volumes.nfs.env.path=%s "+
//			// 代码主机环境目录挂载地址
//			"--set securityContext.runAsUser=%d",
//		// 角色
//		languageEnum, Id, labId,
//		// pod的名字
//		chartName,
//		// chart模板名
//		port,
//		// 端口
//		repository,
//		// 镜像地址
//		tag,
//		// 标签
//		mountEnvDocker,
//		// 容器内环境挂载目录
//		g.Cfg().GetString("ide.storage.host"),
//		// 代码储存主机
//		mountWorkSpaceRemote,
//		// 代码主机工作目录挂载地址
//		g.Cfg().GetString("ide.storage.host"),
//		// 代码储存主机
//		mountEnvRemote,
//		// 代码主机环境目录挂载地址
//		runAsUser,
//		// 角色
//	)
//	// export KUBECONFIG=/etc/rancher/k3s/k3s.yaml 这个必须在以session为单位的执行
//	cmd = "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml && " + cmd
//	if _, err = utils.DeploymentSsh.ExecCmd(cmd); err != nil {
//		return 0, err
//	}
//	return port, nil
//}
//
//func (receiver *k3sIDEService) execUninstallTheiaDeployment(languageEnum int, userId int, labId int) (err error) {
//	// 删除容器
//	cmd := fmt.Sprintf("helm uninstall myIde-%d-%d-%d", languageEnum, userId, labId)
//	cmd = "export KUBECONFIG=/etc/rancher/k3s/k3s.yaml && " + cmd
//	if _, err = utils.DeploymentSsh.ExecCmd(cmd); err != nil {
//		return err
//	}
//	return nil
//}
