package main

import (
	"fmt"
	"monitor/config"
	logger "monitor/logger"
	"monitor/route"
)

func main() {
	// 初始化日志
	if err := logger.Setup(config.Get().Logger); err != nil {
		fmt.Println(err)
		return
	}

	r := route.Route()

	err := r.Run(fmt.Sprintf("%s:%d", config.Get().Server.Host, config.Get().Server.Port))
	if err != nil {
		fmt.Println(err)
	}
}
