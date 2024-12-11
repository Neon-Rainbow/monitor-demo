package main

import (
	"context"
	"fmt"
	"monitor/config"
	"monitor/consul"
	"monitor/etcd"
	logger "monitor/logger"
	"monitor/metrics"
	"monitor/route"

	"go.uber.org/zap"
)

func main() {
	// 初始化日志
	if err := logger.Setup(config.Get().Logger); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化指标
	metrics.NewMetrics().AutoUpdateMetrics()

	// 初始化路由
	r := route.Route()

	// 注册 etcd 服务
	err := etcd.GetService().Register()
	if err != nil {
		zap.L().Fatal("注册服务失败", zap.Error(err))
	}

	// 从 etcd 中获取服务列表
	etcd.AutoFetchServices(context.Background())

	// 注册 consul 服务
	err = consul.GetService().Register()
	if err != nil {
		zap.L().Fatal("注册服务失败", zap.Error(err))
	}

	err = r.Run(fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port))
	if err != nil {
		fmt.Println(err)
	}
}
