package middleware

import (
	"net/http"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var rateLimitIncrExpireScript = redis.NewScript(`
local current = redis.call("INCR", KEYS[1])
if current == 1 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
else
  local ttl = redis.call("PTTL", KEYS[1])
  if ttl < 0 then
    redis.call("PEXPIRE", KEYS[1], ARGV[1])
  end
end
return current
`)

func RateLimit(limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		key := common.RedisKeyRateLimit.F(c.FullPath(), ip)
		ctx := c.Request.Context()
		windowMs := window.Milliseconds()
		if windowMs <= 0 {
			windowMs = 1
		}
		cnt, err := rateLimitIncrExpireScript.Run(ctx, cache.RDB, []string{key}, windowMs).Int64()
		if err == nil && cnt > limit {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": ""})
			return
		}
		c.Next()
	}
}
