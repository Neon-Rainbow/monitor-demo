package route

import (
	"monitor/config"
	"monitor/domain"
	"monitor/logger"
	"monitor/middleware"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

var route *gin.Engine

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

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
		zap.L().Info("health", zap.String("message", "ok"))
	})

	route = r

	return r
}

// ExtractRoutes 用于提取路由信息
// 该函数会提取所有路由信息, 并返回一个 Route 数组
//
// 参数:
//   - router: *gin.Engine: gin.Engine 指针
//
// 返回值:
//   - []domain.Route: Route 数组
func ExtractRoutes() []domain.Route {
	var routes []domain.Route

	for _, r := range route.Routes() {
		routes = append(routes, domain.Route{
			Path:   r.Path,
			Method: r.Method,
		})
	}

	return routes
}
