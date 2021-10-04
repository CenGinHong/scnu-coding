package utils

import (
	"github.com/desops/sshpool"
	"github.com/gogf/gf/frame/g"
	"golang.org/x/crypto/ssh"
	"time"
)

// @Author: 陈健航
// @Date: 2021/4/20 19:23
// @Description:

var DeploymentSsh = newDeploymentSshPoolComponent()

type deploymentSshComponent struct {
	*sshpool.Pool
}

func newDeploymentSshPoolComponent() (d deploymentSshComponent) {
	var config *ssh.ClientConfig
	switch g.Cfg().GetString("ide.deploymentType") {
	// 使用docker部署
	case "docker":
		config = &ssh.ClientConfig{
			User:            g.Cfg().GetString("ide.deployment.docker.user"),
			Auth:            []ssh.AuthMethod{ssh.Password(g.Cfg().GetString("ide.deployment.docker.pass"))},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
		// 使用k3s部署
	case "k3s":
		config = &ssh.ClientConfig{
			User:            g.Cfg().GetString("ide.deployment.k3s.user"),
			Auth:            []ssh.AuthMethod{ssh.Password(g.Cfg().GetString("ide.deployment.k3s.pass"))},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         30 * time.Second,
		}
	default:
		panic("不支持的部署方式")
	}
	// ssh client 初始化
	pool := sshpool.New(config, nil)
	d = deploymentSshComponent{pool}
	session, err := d.Get(g.Cfg().GetString("ide.deployment.docker.host"))
	if err != nil {
		panic("session连接失败")
	}
	defer session.Put()
	// 测试下
	if err = session.Run("echo test"); err != nil {
		panic(err)
	}
	return d
}

func (receiver *deploymentSshComponent) ExecCmd(cmd string) (output string, err error) {
	var session *sshpool.Session
	// 申请session，因为连接限制，所以需要不断切换协程循环申请
	host := g.Cfg().GetString("ide.deployment.docker.host")
	for session == nil {
		session, err = receiver.Get(host)
		if err != nil {
			// 当前session已申请满，切换让出当前协程
			if err.Error() == "ssh: rejected: administratively prohibited (open failed)" {
				time.Sleep(1 * time.Nanosecond)
			} else {
				return "", err
			}
		}
	}
	// 放回池子
	defer session.Put()
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err = session.RequestPty("xterm", 25, 80, modes); err != nil {
		return "", err
	}
	bytes, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	output = string(bytes)
	return output, nil
}
