package middleware

import (
	"context"
	"net/http"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"

	"github.com/gin-gonic/gin"
)

func RateLimit(limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := common.RedisKeyRateLimit.F(c.FullPath(), ip)
		cnt, err := cache.RDB.Incr(context.Background(), key).Result()
		if err == nil && cnt == 1 {
			_ = cache.RDB.Expire(context.Background(), key, window).Err()
		}
		if err == nil && cnt > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limited"})
			return
		}
		c.Next()
	}
}
