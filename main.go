package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"monitor/config"
	"monitor/etcd"
	logger "monitor/logger"
	"monitor/metrics"
	"monitor/route"
)

func main() {
	// 初始化日志
	if err := logger.Setup(config.Get().Logger); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化指标
	metrics.NewMetrics().AutoUpdateMetrics()

	_ = etcd.GetClient()

	// 初始化路由
	r := route.Route()

	// 注册服务
	err := etcd.GetService().Register()
	if err != nil {
		zap.L().Fatal("注册服务失败", zap.Error(err))
	}

	// 从 etcd 中获取服务列表
	etcd.AutoFetchServices(context.Background())

	err = r.Run(fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port))
	if err != nil {
		fmt.Println(err)
	}
}
