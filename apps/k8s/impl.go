package k8s

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// Exec 在 Kubernetes 的一个 Pod 上执行命令
// cmd 是要在 Pod 上执行的命令
// podInfo 包含 Pod 的基本信息
func (k *K8s) Exec(cmd []string, podInfo *PodInfo) error {
	// 创建执行请求
	req, err := k.createExecRequest(cmd, podInfo)
	if err != nil {
		return fmt.Errorf("创建执行请求失败：%v", err)
	}

	// 创建执行器
	exec, err := k.createExecutor(req)
	if err != nil {
		return fmt.Errorf("创建执行器失败：%v", err)
	}

	// 执行命令并获取输出
	stdout, stderr, err := k.streamExec(exec)
	if err != nil {
		zap.L().Error("执行命令失败", zap.String("stdout", stdout), zap.String("stderr", stderr), zap.String("error", err.Error()))
		return fmt.Errorf("执行命令失败：%v", err)
	}

	zap.L().Info("命令执行成功", zap.String("stdout", stdout), zap.String("stderr", stderr))
	return nil
}

// createExecRequest 创建执行请求
// cmd 是要在 Pod 上执行的命令
// podInfo 包含 Pod 的基本信息
func (k *K8s) createExecRequest(cmd []string, podInfo *PodInfo) (*rest.Request, error) {
	return k.ClientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podInfo.PodName).
		Namespace(podInfo.NameSpace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Stdin:   false,
			Stdout:  true,
			Stderr:  true,
			TTY:     false,
			Command: cmd,
		}, scheme.ParameterCodec), nil
}

// createExecutor 创建执行器
// req 是执行请求
func (k *K8s) createExecutor(req *rest.Request) (remotecommand.Executor, error) {
	return remotecommand.NewSPDYExecutor(k.Config, "POST", req.URL())
}

// streamExec 执行命令并获取输出
// exec 是执行器
func (k *K8s) streamExec(exec remotecommand.Executor) (string, string, error) {
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	err := exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: stdout,
		Stderr: stderr,
		Tty:    false,
	})
	if err != nil {
		return stdout.String(), stderr.String(), err
	}
	return stdout.String(), stderr.String(), nil
}

// GetPod 获取 Pod 的详细信息
func (k *K8s) GetPod(podInfo *PodInfo) (*PodInfo, error) {
	pod, err := k.ClientSet.CoreV1().Pods(podInfo.NameSpace).Get(context.TODO(), podInfo.PodName, metav1.GetOptions{})
	if err != nil {
		zap.L().Error("GetPod failed", zap.String("message", "根据pod名称获取podIP失败"), zap.String("error", err.Error()))
		return nil, err
	}
	podInfo.PodIP = pod.Status.PodIP
	zap.L().Info("GetPod succeeded", zap.String("PodIP", podInfo.PodIP))
	return podInfo, nil
}

// ListEvents 获取 Pod 事件
func (k *K8s) ListEvents(podInfo *PodInfo) (*PodInfo, error) {
	podEvents, err := k.ClientSet.CoreV1().Events(podInfo.NameSpace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		zap.L().Error("ListEvents failed", zap.String("message", "获取pod事件失败"), zap.String("error", err.Error()))
		return podInfo, err
	}

	for _, v := range podEvents.Items {
		if v.Type == "Normal" {
			continue
		}
		if strings.HasPrefix(v.Name, podInfo.PodName) {
			podInfo.Events = append(podInfo.Events, v.Message)
		}
	}
	return podInfo, nil
}
