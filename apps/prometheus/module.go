package prometheus

// AlertManager 表示一个告警管理器，其中包含接收者、状态和一组告警
type AlertManager struct {
	Receiver string  `json:"receiver"` // 接收者
	Status   string  `json:"status"`   // 状态
	Alerts   []Alert `json:"alerts"`   // 告警列表
}

// NewAlertManager 创建一个新的 AlertManager 实例
func NewAlertManager() *AlertManager {
	return &AlertManager{
		Alerts: []Alert{},
	}
}

// Alert 表示具体的告警内容
type Alert struct {
	Status       string                 `json:"status"`       // 告警状态
	Labels       map[string]interface{} `json:"labels"`       // 标签
	Annotations  Annotations            `json:"annotations"`  // 注释
	StartsAt     string                 `json:"startsAt"`     // 告警开始时间
	EndsAt       string                 `json:"endsAt"`       // 告警结束时间
	GeneratorUrl string                 `json:"generatorURL"` // 生成器 URL
	FingerPrint  string                 `json:"fingerprint"`  // 指纹
	DurationAt   string                 `json:"durationAt"`   // 持续时间
	Events       []string               // 事件列表
}

// NewAlert 初始化 Alert 实例
func NewAlert() *Alert {
	return &Alert{
		Labels:      make(map[string]interface{}),
		Annotations: Annotations{},
		Events:      []string{},
	}
}

// Annotations 告警注释信息
type Annotations struct {
	Description  string `json:"description"`  // 描述
	Summary      string `json:"summary"`      // 概要
	CurrentValue string `json:"currentvalue"` // 当前值 此标签为自定义标签
}

// DumpData 解析使用
type DumpData struct {
	PodName   string `json:"podName"`   // Pod 名称
	PodIP     string `json:"podIP"`     // Pod IP
	Title     string `json:"title"`     // 标题
	CreatedAt string `json:"createdAt"` // 创建时间
	Url       string `json:"url"`       // URL
}

// NewDumpData 初始化 DumpData 实例
func NewDumpData() *DumpData {
	return &DumpData{}
}
