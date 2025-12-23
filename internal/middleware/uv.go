package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func dailyUVIdentity(c *gin.Context) string {
	if uid, ok := UID(c); ok && uid > 0 {
		return "u:" + itoa(int(uid))
	}
	return "i:" + c.ClientIP()
}

func DailyUV(rdb *redis.Client, exclude ...string) gin.HandlerFunc {
	ex := make(map[string]struct{}, len(exclude))
	for _, p := range exclude {
		ex[p] = struct{}{}
	}
	return func(c *gin.Context) {
		if _, ok := ex[c.FullPath()]; ok {
			c.Next()
			return
		}
		now := time.Now()
		key := "uv:" + now.Format("2006-01-02")
		id := dailyUVIdentity(c)
		added, err := rdb.SAdd(context.Background(), key, id).Result()
		if err == nil && added > 0 {
			ttl, _ := rdb.TTL(context.Background(), key).Result()
			if ttl <= 0 {
				_ = rdb.Expire(context.Background(), key, 60*24*time.Hour).Err()
			}
		}
		c.Next()
	}
}
