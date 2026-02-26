package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"

	"github.com/gin-gonic/gin"
)

func RateLimitUser(limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uid string
		if v, ok := c.Get("uid"); ok {
			uid = toStringUint(v)
		} else {
			uid = c.ClientIP()
		}
		key := common.RedisKeyRateLimitUser.F(c.FullPath(), uid)
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

func toStringUint(v interface{}) string {
	switch t := v.(type) {
	case uint:
		return itoa(int(t))
	case int:
		return itoa(t)
	default:
		return "0"
	}
}
func itoa(i int) string { return fmt.Sprintf("%d", i) }
