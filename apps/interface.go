package apps

import (
	"github.com/tchuaxiaohua/promAlerter/apps/k8s"
	"github.com/tchuaxiaohua/promAlerter/apps/notify"
	"github.com/tchuaxiaohua/promAlerter/apps/prometheus"
	"github.com/tchuaxiaohua/promAlerter/svc"
)

type K8sService interface {
	Exec(cmd []string, pod *k8s.PodInfo) error
	GetPod(pod *k8s.PodInfo) (*k8s.PodInfo, error)
	ListEvents(pod *k8s.PodInfo) (*k8s.PodInfo, error)
}

// NotifierService 告警发送接口
type NotifierService interface {
	ProcessAlert(alert prometheus.Alert, svc *svc.AppService, notificationConfig *notify.NotificationConfig) error
}
