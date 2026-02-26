package middleware

import (
	"strings"

	"cvforge/internal/common"
	conf "cvforge/internal/module/config"

	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	sys := conf.NewService()
	origins := sys.Get(string(common.ConfigKeyCORSOrigins))
	allowed := strings.Split(origins, ",")
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin != "" {
			for _, o := range allowed {
				o = strings.TrimSpace(o)
				if o == "*" || origin == o {
					c.Header("Access-Control-Allow-Origin", origin)
					c.Header("Access-Control-Allow-Credentials", "true")
					c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
					c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
					break
				}
			}
		}
		if c.Request.Method == "OPTIONS" {
			c.Status(204)
			c.Abort()
			return
		}
		c.Next()
	}
}
