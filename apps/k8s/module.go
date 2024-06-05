package k8s

import (
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"runtime"
	"sync"
	"time"
)

// K8s 结构体
type K8s struct {
	ClientSet kubernetes.Interface
	Config    *rest.Config
}

type PodInfo struct {
	PodName      string
	NameSpace    string
	PodIP        string
	RestartCount int
	Events       []string
}

func NewPodInfo() *PodInfo {
	return &PodInfo{}
}

// NewK8s 获取k8s实例
// PodName pod名称 NameSpace 命名空间 PodIP pod ip 三个参数需要再告警的时候补全
func NewK8s() *K8s {
	var err error
	var k8sConfig *rest.Config

	if runtime.GOOS == "windows" {
		k8sConfig, err = clientcmd.BuildConfigFromFlags("", "D:\\工作\\goapps\\promDingTalk\\etc\\config")
	} else {
		k8sConfig, err = rest.InClusterConfig()
	}
	if err != nil {
		zap.L().Error("config failed", zap.String("error", err.Error()))
		return nil
	}
	// 初始化客户端
	k8sClient, err := kubernetes.NewForConfig(k8sConfig)
	if err != nil {
		zap.L().Error("client init failed", zap.String("error", err.Error()))
		return nil
	}
	return &K8s{
		ClientSet: k8sClient,
		Config:    k8sConfig,
	}
}

type TokenBucket struct {
	mu         sync.Mutex
	capacity   int           // 桶的容量
	tokens     int           // 当前令牌数
	refillRate time.Duration // 令牌补充速率，例如每秒补充一个令牌
	lastRefill time.Time     // 上次补充令牌的时间
}

func NewTokenBucket(capacity int, refillRate time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity, // 初始化时桶满
		refillRate: refillRate,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)
	refillAmount := elapsed / tb.refillRate
	newTokens := tb.tokens + int(refillAmount.Seconds())

	if newTokens > tb.capacity {
		newTokens = tb.capacity
	}

	tb.tokens = newTokens
	tb.lastRefill = now

	if tb.tokens < 1 {
		return false // 没有足够的令牌，不能继续
	}

	tb.tokens--
	return true // 有足够的令牌，可以继续
}
