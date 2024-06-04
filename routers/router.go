package routers

import (
	"github.com/tchuaxiaohua/promDingTalk/config"
	"github.com/tchuaxiaohua/promDingTalk/routers/notify"
	"github.com/tchuaxiaohua/promDingTalk/svc"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
)

// Router 项目路由
func Router(r gin.IRouter, appCtx *svc.AppService) {
	r.POST("/api/notify/:app", notify.Notify(appCtx))
}

func InitRouter(appCtx *svc.AppService) *gin.Engine {
	r := gin.New()
	//gin 接入自定义logger
	r.Use(ginzap.Ginzap(config.Log, time.RFC3339, false), gin.Recovery())
	// 用户路由 加载
	Router(r, appCtx)
	return r
}
