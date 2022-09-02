package middlewaresGin

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start)
		statusCode := c.Writer.Status()
		msg := []zap.Field{
			//日志类型
			zap.String("type", "webase-server-request"),
			//请求用户的 IP
			zap.String("ip", c.ClientIP()),
			//请求的 RequestURI
			zap.String("uri", c.Request.RequestURI),
			//请求的方法
			zap.String("method", c.Request.Method),
			//http状态码
			zap.Int("statusCode", statusCode),
			//请求花费时间
			zap.Duration("cost", cost),
		}
		if statusCode > 499 {
			zap.L().Error("请求响应", msg...)
		} else if statusCode > 399 {
			zap.L().Warn("请求响应", msg...)
		} else {
			zap.L().Info("请求响应", msg...)
		}
	}
}
