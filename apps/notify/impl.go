package notify

import (
	"fmt"

	"github.com/tchuaxiaohua/promDingTalk/apps/prometheus"
	"github.com/tchuaxiaohua/promDingTalk/config"
	"github.com/tchuaxiaohua/promDingTalk/svc"
)

func (n *NotificationConfig) ProcessAlert(alert prometheus.Alert, svc *svc.AppService, notificationConfig *NotificationConfig) error {
	// 计算告警持续时间
	alert.DurationTime()
	// 时间格式化解析处理
	if err := alert.TimeFormat(); err != nil {
		return fmt.Errorf("告警时间解析失败: %w", err)
	}

	// 判断是否需要执行 dump 操作
	if shouldDump(alert, svc.Config.Jvm.IsDump) {
		if err := alert.Dump(svc.K8s); err != nil {
			return fmt.Errorf("dump 操作失败: %w", err)
		}
	}

	// 处理 pod 事件
	if _, ok := alert.Labels["pod"]; ok {
		alert.GetEvents(svc.K8s)
	}

	// 发送告警
	notifiers, err := getNotifier(&svc.Config.Channels, notificationConfig.AppName)
	if err != nil {
		return fmt.Errorf("初始化告警器失败: %w", err)
	}
	// 循环初始化告警器
	for _, notifier := range notifiers {
		if err := notifier.Send(&alert, notificationConfig); err != nil {
			return fmt.Errorf("发送告警失败: %w", err)
		}
	}
	return nil
}

// ShouldDump 判断是否需要执行dump操作
func shouldDump(alert prometheus.Alert, isDumpEnabled bool) bool {
	return alert.Labels["jvm_dump"] == "true" && isDumpEnabled
}

// getNotifier 封装告警通知服务
func getNotifier(configs *config.ChannelConfigs, appName string) ([]notifierService, error) {
	var notifiers []notifierService

	// 钉钉
	if configs.DingTalk.Enabled {
		for _, cfg := range configs.DingTalk.Configs {
			if cfg.Name == appName {
				notifiers = append(notifiers, &DingTalkNotifier{Token: cfg.Token, Secret: cfg.Secret})
			}
		}
	}
	// 飞书
	if configs.FeiShu.Enabled {
		for _, cfg := range configs.FeiShu.Configs {
			if cfg.Name == appName {
				notifiers = append(notifiers, &FeiShuNotifier{Token: cfg.Token, Secret: cfg.Secret})
			}
		}
	}
	// 邮件
	if configs.Email.Enabled {
		for _, cfg := range configs.Email.Configs {
			notifiers = append(notifiers, &EmailNotifier{
				SMTPServer:  cfg.SMTPServer,
				SMTPPort:    cfg.SMTPPort,
				Username:    cfg.Username,
				Password:    cfg.Password,
				FromAddress: cfg.FromAddress,
				ToAddress:   cfg.ToAddress,
			})
		}
	}

	if len(notifiers) == 0 {
		return nil, fmt.Errorf("no enabled notification channel found for app: %s", appName)
	}

	return notifiers, nil
}
