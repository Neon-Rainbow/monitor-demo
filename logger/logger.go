package logger

import (
	"bytes"
	"io"
	"monitor/config"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"go.uber.org/zap/zapcore"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// newLogger 创建并初始化zap日志库
func newLogger(cfg *config.Logger) (*zap.Logger, error) {
	// 创建日志文件
	logFile, err := os.Create(cfg.OutputPath)
	if err != nil {
		return nil, err
	}
	// 保证日志文件在程序退出时关闭
	writeSyncer := zapcore.AddSync(logFile)

	// 设置日志级别,从配置文件中读取
	var level zapcore.Level

	logLevelMapper := map[int]zapcore.Level{
		-1: zap.DebugLevel,
		0:  zap.InfoLevel,
		1:  zap.WarnLevel,
		2:  zap.ErrorLevel,
		3:  zap.DPanicLevel,
		4:  zap.PanicLevel,
		5:  zap.FatalLevel,
	}

	level = zap.InfoLevel
	if l, ok := logLevelMapper[cfg.Level]; ok {
		level = l
	}

	// 设置日志格式（JSON 或 Console）
	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	} else {
		encoder = zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	}

	// 设置日志输出路径和错误输出路径
	cores := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout), // 输出到标准输出
			writeSyncer,                // 输出到文件
		),
		level,
	)

	// 构建zap日志对象
	logger := zap.New(cores, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return logger, nil
}

// GinLogger 接收gin框架默认的日志
func GinLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 由于c.Request.Body只允许读一次,读完之后内容会被删除,因此在这里我们读取Body中的内容,然后存储在变量中
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		c.Next()

		cost := time.Since(start)
		logger.Info(path,
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.Any("head", c.Request.Header),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("body", string(bodyBytes)),
			zap.String("ip", c.ClientIP()),
			zap.String("user-agent", c.Request.UserAgent()),
			zap.String("errors", c.Errors.ByType(gin.ErrorTypePrivate).String()),
			zap.Duration("cost", cost),
			zap.String("request_id", c.Request.Header.Get("X-Request-Id")),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(logger *zap.Logger, stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					logger.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					logger.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
