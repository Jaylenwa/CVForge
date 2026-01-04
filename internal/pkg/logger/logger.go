package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var root *zap.Logger

func Init() {
	if root == nil {
		root, _ = zap.NewProduction()
	}
}

func L() *zap.Logger {
	if root == nil {
		Init()
	}
	return root
}

func WithCtx(c *gin.Context) *zap.Logger {
	l := L()
	if c == nil {
		return l
	}
	rid := c.Writer.Header().Get("X-Request-ID")
	return l.With(
		zap.String("request_id", rid),
		zap.String("method", c.Request.Method),
		zap.String("path", c.FullPath()),
	)
}
