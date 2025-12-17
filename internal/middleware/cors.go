package middleware

import (
	"strings"

	"openresume/internal/infra/config"

	"github.com/gin-gonic/gin"
)

func CORS(cfg config.Config) gin.HandlerFunc {
	allowed := strings.Split(cfg.CORSOrigins, ",")
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
