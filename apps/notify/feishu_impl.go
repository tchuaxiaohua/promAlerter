package notify

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/tchuaxiaohua/promDingTalk/apps/prometheus"
)

// Send 通过飞书发送告警
func (f *FeiShuNotifier) Send(alert *prometheus.Alert, notificationConfig *NotificationConfig) error {
	webhook := fmt.Sprintf("%s&sign=%s", f.Token, f.Secret)
	message := fmt.Sprintf(`{"msg_type": "text", "content": {"text": "%s"}}`, notificationConfig.Title)
	req, err := http.NewRequest("POST", webhook, bytes.NewBuffer([]byte(message)))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send feishu notification, status code: %d", resp.StatusCode)
	}
	return nil
}
