package notify

import (
	"github.com/tchuaxiaohua/promAlerter/apps/notify"
	"github.com/tchuaxiaohua/promAlerter/apps/prometheus"
	"github.com/tchuaxiaohua/promAlerter/svc"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Notify 路由处理函数
func Notify(svc *svc.AppService) gin.HandlerFunc {
	return func(c *gin.Context) {
		//	接收alertManager告警消息
		alertObj := prometheus.NewAlertManager()
		if err := c.ShouldBindJSON(alertObj); err != nil {
			zap.L().Error("Should Bind error", zap.String("error", err.Error()), zap.String("message", "参数解释错误"))
			return
		}
		// 获取告警渠道 这里用来判断符合的机器人
		appName := c.Param("app")
		// 初始化告警逻辑处理实例对象
		notifier := notify.NewNotificationConfig(appName)
		// 解析告警信息
		for _, alert := range alertObj.Alerts {
			if err := notifier.ProcessAlert(alert, svc, notifier); err != nil {
				zap.L().Error("处理告警失败", zap.String("error", err.Error()))
				continue
			}
		}

		c.JSON(200, gin.H{
			"code": 200,
			"msg":  "ok",
			"data": alertObj,
		})
	}
}
