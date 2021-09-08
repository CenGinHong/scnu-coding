package utils

// @Author: 陈健航
// @Date: 2021/1/4 22:41
// @Description:

import (
	"github.com/gogf/gf/frame/g"
	"gopkg.in/gomail.v2"
)

var MailUtil = newMailUtil()

type mailUtil struct {
	dialer *gomail.Dialer
}

func newMailUtil() (m mailUtil) {
	m = mailUtil{}
	//定义邮箱服务器连接信息，如果是网易邮箱 pass填密码，qq邮箱填授权码
	user := g.Cfg().GetString("mail.user")
	pass := g.Cfg().GetString("mail.pass")
	host := g.Cfg().GetString("mail.host")
	port := g.Cfg().GetInt("mail.port")
	m.dialer = gomail.NewDialer(host, port, user, pass)
	return m
}

func (u *mailUtil) SendMail(mailTo string, subject string, body string) (err error) {
	m := gomail.NewMessage()
	user := g.Cfg().GetString("mail.user")
	m.SetHeader("From", m.FormatAddress(user, "在线编程系统")) //这种方式可以添加别名，即“XX官方”
	//说明：如果是用网易邮箱账号发送，以下方法别名可以是中文，如果是qq企业邮箱，以下方法用中文别名，会报错，需要用上面此方法转码
	//m.SetHeader("From", "FB Sample"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“FB Sample”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	//m.SetHeader("From", mailConn["user"])
	m.SetHeader("To", mailTo)       //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文
	if err = u.dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
