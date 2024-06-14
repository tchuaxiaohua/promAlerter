package notify

import (
	"github.com/tchuaxiaohua/promAlerter/apps/prometheus"
	"github.com/tchuaxiaohua/promAlerter/utils"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

func (e *EmailNotifier) Send(alert *prometheus.Alert, notificationConfig *NotificationConfig) error {

	//	解析模板渲染告警信息
	data, err := utils.ParseTemplate("templates/email.html", alert)
	if err != nil {
		zap.L().Error("邮件模板解析失败", zap.String("error", err.Error()), zap.String("message", "模板解析失败"))
	}

	m := gomail.NewMessage()

	// 设置邮件发送者、接收者、主题和正文
	m.SetHeader("From", e.FromAddress)
	m.SetHeader("To", e.ToAddress)
	m.SetHeader("Subject", "Prometheus告警")
	m.SetBody("text/html", data)

	// 设置SMTP服务器信息
	d := gomail.NewDialer(e.SMTPServer, e.SMTPPort, e.FromAddress, e.Password)

	// 连接SMTP服务器并发送邮件
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
