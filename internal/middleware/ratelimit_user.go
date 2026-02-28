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
local ttl = redis.call("PTTL", KEYS[1])
if current == 1 or ttl < 0 then
  redis.call("PEXPIRE", KEYS[1], ARGV[1])
  ttl = redis.call("PTTL", KEYS[1])
end
return {current, ttl}
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
		out, err := rateLimitUserIncrExpireScript.Run(ctx, cache.RDB, []string{key}, windowMs).Slice()
		if err != nil || len(out) < 2 {
			c.Next()
			return
		}
		cnt, okCnt := asInt64(out[0])
		ttlMs, okTTL := asInt64(out[1])
		if !okCnt || !okTTL {
			c.Next()
			return
		}
		if cnt > limit {
			retryAfterSeconds := int64(window.Seconds())
			if retryAfterSeconds <= 0 {
				retryAfterSeconds = 1
			}
			if ttlMs > 0 {
				retryAfterSeconds = (ttlMs + 999) / 1000
			}
			remaining := limit - cnt
			if remaining < 0 {
				remaining = 0
			}
			c.Header("Retry-After", fmt.Sprintf("%d", retryAfterSeconds))
			c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
			c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limited", "code": "rate_limited", "retryAfterSeconds": retryAfterSeconds})
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
