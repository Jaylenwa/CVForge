package middleware

import (
	"context"
	"net/http"
	"time"

	"openresume/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func RateLimit(rdb *redis.Client, limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := common.RedisKeyRateLimit.F(c.FullPath(), ip)
		cnt, err := rdb.Incr(context.Background(), key).Result()
		if err == nil && cnt == 1 {
			_ = rdb.Expire(context.Background(), key, window).Err()
		}
		if err == nil && cnt > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limited"})
			return
		}
		c.Next()
	}
}
