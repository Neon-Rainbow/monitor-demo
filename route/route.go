package route

import (
	"monitor/config"
	"monitor/logger"
	"monitor/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Route 用于定义路由
func Route() *gin.Engine {
	r := gin.New()

	r.Use(logger.GinLogger(zap.L()), logger.GinRecovery(zap.L(), true))

	switch config.Get().Server.Mode {
	case "release":
		gin.SetMode(gin.ReleaseMode)
	case "debug":
		gin.SetMode(gin.DebugMode)
	case "test":
		gin.SetMode(gin.TestMode)
	default:
		gin.SetMode(gin.DebugMode)
	}

	// 注册路由
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 注册中间件
	r.Use(middleware.PrometheusMiddleware())

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
		zap.L().Info("ping", zap.String("message", "pong"))
	})

	return r
}
