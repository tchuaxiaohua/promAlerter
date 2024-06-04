package prometheus

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/tchuaxiaohua/promDingTalk/apps/k8s"
	"github.com/tchuaxiaohua/promDingTalk/config"
	"github.com/tchuaxiaohua/promDingTalk/utils"
	"github.com/tchuaxiaohua/promDingTalk/vars"
	"go.uber.org/zap"
	"time"
)

// Dump 执行 JVM 内存快照导出操作
func (a *Alert) Dump(k8sClient *k8s.K8s) error {
	podObj := k8s.NewPodInfo()
	podObj.PodName = a.getMap("pod")
	podObj.NameSpace = a.getMap("namespace")

	if a.Status == "resolved" {
		return errors.New(fmt.Sprintf("【%s】告警恢复，无需处理", podObj.PodName))
	}

	currentValue := a.Annotations.CurrentValue
	// 判断是否超过最大时间限制
	if duration := time.Since(utils.PtrTime(a.StartsAt)); duration > time.Duration(config.C().Jvm.DumpTsMax)*time.Hour {
		return errors.New(fmt.Sprintf("【%s】告警时间超过%dH,已执行dump操作,告警触发时间:%s", podObj.PodName, config.C().Jvm.DumpTsMax, a.StartsAt))
	}
	// 判断是否超过最小时间限制
	if duration := time.Since(utils.PtrTime(a.StartsAt)); !(duration > time.Duration(config.C().Jvm.DumpTsMin)*time.Hour) {
		return errors.New(fmt.Sprintf("【%s】告警触发时间不足%dH,告警触发时间:%s", podObj.PodName, config.C().Jvm.DumpTsMin, a.StartsAt))
	}
	// 判断当前内存使用率是否在指定范围内
	if value := utils.PtrInt(currentValue); value < config.C().Jvm.DumpMin || value > config.C().Jvm.DumpMax {
		return errors.New(fmt.Sprintf("【%s】当前值内存使用率已超过%d%%,当前值:%d%%,退出执行dump操作！！！", podObj.PodName, config.C().Jvm.DumpMax, value))
	}
	// 执行 JVM 内存快照导出
	if err := k8sClient.Exec(vars.CmdDump, podObj); err != nil {
		return fmt.Errorf("failed to execute JVM memory dump: %w", err)
	}
	// 上传内存快照文件至 OSS
	if err := k8sClient.Exec(vars.CmdUploadDump, podObj); err != nil {
		return fmt.Errorf("failed to upload memory dump file to OSS: %w", err)
	}
	return nil
}

// getMap 从 Labels 中获取值
func (a *Alert) getMap(key string) string {
	v, ok := a.Labels[key]
	if ok {
		return v.(string)
	}
	return ""
}

func (a *Alert) GetEvents(k8sClient *k8s.K8s) {
	podObj := k8s.NewPodInfo()
	// 取podEvents
	podObj.PodName = a.getMap("pod")
	podObj.NameSpace = a.getMap("namespace")
	k8sPodInfo, err := k8sClient.ListEvents(podObj)
	if err != nil {
		zap.L().Error(fmt.Sprintf("failed to get pod events: %v", err))
	}
	a.Events = k8sPodInfo.Events
}

// TimeFormat 解析告警时间为标准时间
func (a *Alert) TimeFormat() error {
	layout := "2006-01-02 15:04:05"

	// 告警触发时间
	t, err := time.Parse(time.RFC3339, a.StartsAt)
	if err != nil {
		return err
	}
	a.StartsAt = t.In(time.Local).Format(layout)

	// 告警结束时间
	if len(a.EndsAt) > 0 {
		t1, err := time.Parse(time.RFC3339, a.EndsAt)
		if err != nil {
			return err
		}
		a.EndsAt = t1.In(time.Local).Format(layout)
	}
	return nil
}

// DurationTime 计算告警持续时间
func (a *Alert) DurationTime() {
	// 解析开始时间 转化为 time.Time 类型
	startTime, err := time.Parse(time.RFC3339, a.StartsAt)
	if err != nil {
		zap.L().Error("时间解析错误", zap.String("err", err.Error()))
	}
	// 处理结束时间: 如果告警恢复 则使用结束时间 否则使用当前时间
	var endTime time.Time
	if a.Status == "resolved" {
		endTime, err = time.Parse(time.RFC3339, a.EndsAt)
		if err != nil {
			zap.L().Error("时间解析错误", zap.String("err", err.Error()))
		}
	} else {
		endTime = time.Now()
	}
	// 计算时间差
	duration := endTime.Sub(startTime)
	// 格式化持续时间，只显示天、小时和分钟
	a.DurationAt = formatDuration(duration)
}

// formatDuration 格式化持续时间，只显示天、小时和分钟
func formatDuration(duration time.Duration) string {
	// 转换持续时间为秒 并计算搞持续时间 天、小时和分钟数
	seconds := int(duration.Seconds())
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24

	// 计算剩余小时数（超过整天的部分）
	remainingHours := hours % 24
	// 计算剩余分钟数（超过整小时的部分）
	remainingMinutes := minutes % 60
	// 拼接结果
	var result string
	if days > 0 {
		result += fmt.Sprintf("%d天", days)
	}
	if remainingHours > 0 {
		result += fmt.Sprintf("%d小时", remainingHours)
	}
	if remainingMinutes > 0 {
		result += fmt.Sprintf("%d分钟", remainingMinutes)
	}

	if result == "" {
		result = "不到1分钟"
	}

	return result
}
