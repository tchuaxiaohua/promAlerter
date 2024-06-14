package notify

import (
	"github.com/tchuaxiaohua/promAlerter/apps/prometheus"
)

// notifierService 告警发送接口
type notifierService interface {
	Send(alert *prometheus.Alert, notification *NotificationConfig) error
}

// NotificationConfig 告警配置
type NotificationConfig struct {
	Title   string // 告警标题
	AppName string // 告警渠道
}

// NewNotificationConfig 告警配置初始化函数
func NewNotificationConfig(robotName string) *NotificationConfig {
	return &NotificationConfig{
		AppName: robotName,
	}
}

// DingTalkNotifier  钉钉告警发送器
type DingTalkNotifier struct {
	Token  string `yaml:"token"`
	Secret string `yaml:"secret"`
}

// FeiShuNotifier 飞书告警发送器
type FeiShuNotifier struct {
	Token  string `yaml:"token"`
	Secret string `yaml:"secret"`
}

// EmailNotifier 邮件告警发送器
type EmailNotifier struct {
	SMTPServer  string `yaml:"smtpServer"`
	SMTPPort    int    `yaml:"smtpPort"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	FromAddress string `yaml:"fromAddress"`
	ToAddress   string `yaml:"toAddress"`
}
