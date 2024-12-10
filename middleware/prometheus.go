package middleware

import (
	"monitor/metrics"
	"strconv"

	"github.com/gin-gonic/gin"
)

// PrometheusMiddleware 用于记录请求次数
func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		path := c.FullPath()

		c.Next()

		statusCode := c.Writer.Status()

		// 统计请求次数
		metrics.HttpRequest.AddCounter(method, path, strconv.Itoa(statusCode))
	}
}
