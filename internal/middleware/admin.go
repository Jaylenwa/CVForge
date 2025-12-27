package middleware

import (
	"net/http"

	"openresume/internal/common"
	"openresume/internal/infra/database"
	"openresume/internal/models"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...common.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		uidVal, ok := c.Get("uid")
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var u models.User
		if err := database.DB.First(&u, uidVal).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		if !u.IsActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		match := false
		for _, r := range roles {
			if u.Role == r {
				match = true
				break
			}
		}
		if !match {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		c.Next()
	}
}
