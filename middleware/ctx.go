package middleware

import (
	"github.com/gin-gonic/gin"
)

func UID(c *gin.Context) (uint, bool) {
	uidVal, ok := c.Get("uid")
	if !ok {
		return 0, false
	}
	switch v := uidVal.(type) {
	case uint:
		return v, true
	case int:
		return uint(v), true
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
}
