package main

import (
	"fmt"
	"monitor/config"
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

	// 初始化路由
	r := route.Route()

	err := r.Run(fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port))
	if err != nil {
		fmt.Println(err)
	}
}
