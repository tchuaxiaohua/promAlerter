package k8s

import (
	"fmt"
	"testing"
)

func TestK8sExec(t *testing.T) {
	// 创建 K8s 对象
	client := NewK8s()
	if client == nil {
		t.Errorf("ailed to initialize K8s client")
	}
	// 创建测试用的命令
	cmd := []string{"echo", "hello world"}
	pod := &PodInfo{
		PodName:   "nginx-deployment-7c79c4bf97-52gmr",
		NameSpace: "default",
	}
	// 执行 Exec 函数
	err := client.Exec(cmd, pod)
	if err != nil {
		t.Errorf("Exec failed with error: %v", err)
	}
}

func TestK8sGetPod(t *testing.T) {
	// 创建 K8s 对象
	client := NewK8s()
	if client == nil {
		t.Errorf("ailed to initialize K8s client")
	}
	// 创建测试用的命令
	pod := &PodInfo{
		PodName:   "nginx-deployment-7c79c4bf97-52gmr",
		NameSpace: "default",
	}
	// 执行 Exec 函数
	podObj, err := client.GetPod(pod)
	if err != nil {
		t.Errorf("Exec failed with error: %v", err)
	}
	fmt.Println(podObj.PodIP)
	fmt.Println(podObj.RestartCount)
}
