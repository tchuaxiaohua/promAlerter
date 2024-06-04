package main

import (
	"github.com/tchuaxiaohua/promDingTalk/config"
	"github.com/tchuaxiaohua/promDingTalk/routers"
	"github.com/tchuaxiaohua/promDingTalk/svc"
	"go.uber.org/zap"
	"runtime"
)

func main() {
	// 初始化 全局配置文件
	var configPath string
	if runtime.GOOS == "windows" {
		configPath = "etc/app-dev.yaml"
	} else {
		configPath = "etc/app.yaml"
	}
	if err := config.LoadConfig(configPath); err != nil {
		zap.L().Error("全局配置加载失败", zap.String("err", err.Error()))
	}
	// 初始化 log
	config.InitLogger(config.C().Log.Level)
	// 初始化 全局服务
	appCtx := svc.NewAppService(config.C())
	// 初始化 路由
	g := routers.InitRouter(appCtx)
	// 启动
	if err := g.Run(config.C().App.GetHost()); err != nil {
		zap.L().Error("应用启动失败")
	}
}
