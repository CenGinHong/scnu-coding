package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	component2 "scnu-coding/app/utils"
	"time"
)

var Common = commonService{verCodeCache: component2.NewMyCache()}

type commonService struct{ verCodeCache component2.MyCache }

func (c *commonService) SendVerificationCode(_ context.Context, email string) (err error) {
	//生成6位随机数
	verCode := grand.Digits(6)
	// 存入redis，5分钟有效期
	if err = c.verCodeCache.Set(email, verCode, 5*time.Minute); err != nil {
		return err
	}
	// 准备邮件内容
	emailBody := fmt.Sprintf(`<div style="background-color:#ECECEC; padding: 35px;">
  <table cellpadding="0" align="center"
         style="width: 600px; margin: 0px auto; text-align: left; position: relative; border-top-left-radius: 5px; border-top-right-radius: 5px; border-bottom-right-radius: 5px; border-bottom-left-radius: 5px; font-size: 14px; font-family:微软雅黑, 黑体; line-height: 1.5; box-shadow: rgb(153, 153, 153) 0px 0px 5px; border-collapse: collapse; background-position: initial initial; background-repeat: initial initial;background:#fff;">
      <tbody>
      <tr>
          <th valign="middle"
              style="height: 25px; line-height: 25px; padding: 15px 35px; border-bottom-width: 1px; border-bottom-style: solid; border-bottom-color: #42a3d3; background-color: #49bcff; border-top-left-radius: 5px; border-top-right-radius: 5px; border-bottom-right-radius: 0px; border-bottom-left-radius: 0px;">
              <font face="微软雅黑" size="5" style="color: rgb(255, 255, 255); ">注册成功!</font>
          </th>
      </tr>
      <tr>
          <td>
              <div style="padding:25px 35px 40px; background-color:#fff;">
                  <h2 style="margin: 5px 0px; ">
                      <font color="#333333" style="line-height: 20px; ">
                          <font style="line-height: 22px; " size="4">
                              亲爱的同学</font>
                      </font>
                  </h2>
                  <p>首先感谢您使用本站！这是您的验证码：%s<br>
                  <p align="right">%s</p>
                  <div style="width:700px;margin:0 auto;">
                      <div style="padding:10px 10px 0;border-top:1px solid #ccc;color:#747474;margin-bottom:20px;line-height:1.3em;font-size:12px;">
                          <p>此为系统邮件，请勿回复<br>
                              请保管好您的邮箱，避免账号被他人盗用
                          </p>
                          <p>©***</p>
                      </div>
                  </div>
              </div>
          </td>
      </tr>
      </tbody>
  </table>
</div>
`, verCode, gtime.Datetime())
	// 开一个协程发送邮件
	go func() {
		_ = component2.MailUtil.SendMail(email, "SCNU-CODING实验系统邮件验证码", emailBody)
	}()
	return nil
}

func (c *commonService) CheckVerCode(email string, verCode string) (err error) {
	// 取出数据
	v, err := c.verCodeCache.GetVar(email)
	if err != nil {
		return err
	}
	// 不存在验证码
	if v.IsNil() {
		return gerror.NewCode(-1, "验证码错误")
	}
	// 验证码错误
	if verCode != v.String() {
		return gerror.NewCode(-1, "验证码错误")
	}
	return nil
}
