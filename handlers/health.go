package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var start = time.Now()

func RegisterHealthRoutes(r *gin.RouterGroup) {
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Since(start).Seconds()})
	})
}
