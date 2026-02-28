package middleware

import (
	"fmt"
	"net/http"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

var rateLimitUserIncrExpireScript = redis.NewScript(`
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

func RateLimitUser(limit int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var uid string
		if v, ok := c.Get("uid"); ok {
			uid = toStringUint(v)
		} else {
			uid = c.ClientIP()
		}
		key := common.RedisKeyRateLimitUser.F(c.FullPath(), uid)
		ctx := c.Request.Context()
		windowMs := window.Milliseconds()
		if windowMs <= 0 {
			windowMs = 1
		}
		cnt, err := rateLimitUserIncrExpireScript.Run(ctx, cache.RDB, []string{key}, windowMs).Int64()
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
