package util

import (
	"strconv"

	log "github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

func SendRegisterMail(userMail, checkCode string) {
	//定义收件人
	mailTo := []string{
		userMail,
	}

	//邮件主题为"Hello"
	subject := "趣链网注册验证"

	url := "http://funny.link/" + checkCode

	// 邮件正文
	body := "点击完成邮箱验证：" + url

	SendMail(mailTo, subject, body)
}

func SendMail(mailTo []string, subject string, body string) error {
	log.WithFields(log.Fields{
		"mailTo":  mailTo,
		"subject": subject,
		"body":    body,
	}).Debug("SendMail")

	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "funnylink@126.com",
		"pass": "",
		"host": "smtp.126.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", "funnylink"+"<"+mailConn["user"]+">") //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)                              //发送给多个用户
	m.SetHeader("Subject", subject)                           //设置邮件主题
	m.SetBody("text/html", body)                              //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)

	if err != nil {
		log.Errorf("SendMail %v", err)
	}

	return err
}
