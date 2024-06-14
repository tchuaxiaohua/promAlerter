package notify

import (
	"github.com/CatchZeng/dingtalk/pkg/dingtalk"
	"github.com/tchuaxiaohua/promAlerter/apps/prometheus"
	"github.com/tchuaxiaohua/promAlerter/utils"

	"go.uber.org/zap"
)

// Send 通过钉钉发送告警
func (d *DingTalkNotifier) Send(alert *prometheus.Alert, notificationConfig *NotificationConfig) error {
	//	解析模板渲染告警信息
	data, err := utils.ParseTemplate("templates/notify.html", alert)
	if err != nil {
		zap.L().Error("模板解析失败", zap.String("error", err.Error()), zap.String("message", "模板解析失败"))
	}
	// 实现发送钉钉消息的逻辑
	dingClient := dingtalk.NewClient(d.Token, d.Secret)
	msg := dingtalk.NewMarkdownMessage().SetMarkdown("Prometheus系统告警", data)
	//msg := dingtalk.NewActionCardMessage().SetIndependentJump("Prometheus系统告警", data, nil, "", "")
	_, _, err = dingClient.Send(msg)
	return err
}

//func (d *DingTalkConfig) Send(notification NotificationConfig) error {
//	webhook := fmt.Sprintf("%s&sign=%s", d.Token, d.Secret)
//	message := fmt.Sprintf(`{"msgtype": "text", "text": {"content": "%s"}}`, notification.Title)
//	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer([]byte(message)))
//	if err != nil {
//		return err
//	}
//	req.Header.Set("Content-Type", "application/json")
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		return err
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed to send dingtalk notification, status code: %d", resp.StatusCode)
//	}
//	return nil
//}
