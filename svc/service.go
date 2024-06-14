package svc

import (
	"github.com/tchuaxiaohua/promAlerter/apps/k8s"
	"github.com/tchuaxiaohua/promAlerter/config"
)

type AppService struct {
	Config *config.Config
	K8s    *k8s.K8s
}

func NewAppService(c *config.Config) *AppService {
	return &AppService{
		Config: c,
		K8s:    k8s.NewK8s(),
	}
}
