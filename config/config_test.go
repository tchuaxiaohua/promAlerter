package config

import (
	"fmt"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// 测试用例1：配置文件路径正确
	t.Run("ConfigPathCorrect", func(t *testing.T) {
		// 准备测试数据
		path := "D:\\工作\\goapps\\promDingTalk\\etc\\app-dev.yaml"
		// 调用函数
		err := LoadConfig(path)
		// 断言结果
		if err != nil {
			t.Errorf("LoadConfig(%s) returned error: %v", path, err)
		}

		fmt.Println(111, cfg.Channels.Email)
	})
}
