package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"openresume/internal/pkg/logger"
	"go.uber.org/zap"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		rid := c.Writer.Header().Get("X-Request-ID")
		logger.L().Info("http", zap.String("method", c.Request.Method), zap.String("path", c.FullPath()), zap.Int("status", c.Writer.Status()), zap.Float64("duration", time.Since(start).Seconds()), zap.String("request_id", rid))
	}
}
