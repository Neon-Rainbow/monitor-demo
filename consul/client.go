package consul

import (
	"monitor/config"
	"sync"

	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
)

var (
	consulClient *api.Client
	consulOnce   sync.Once
)

// newConsul 用于创建 Consul 客户端
func newConsul(address string) {
	cfg := api.DefaultConfig()
	cfg.Address = address
	var err error
	consulClient, err = api.NewClient(cfg)
	if err != nil {
		zap.L().Fatal("创建 Consul 客户端失败", zap.Error(err))
		return
	}
}

// GetClient 用于获取 Consul 客户端
func GetClient() *api.Client {
	consulOnce.Do(func() {
		newConsul(config.Get().Consul.Address)
	})
	return consulClient
}


