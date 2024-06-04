package k8s

import (
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"runtime"
)

// K8s 结构体
type K8s struct {
	ClientSet kubernetes.Interface
	Config    *rest.Config
}

type PodInfo struct {
	PodName   string
	NameSpace string
	PodIP     string
	Events    []string
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
