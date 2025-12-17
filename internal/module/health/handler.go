package health

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

var start = time.Now()

type Handler struct{}

func (h *Handler) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "uptime": time.Since(start).Seconds()})
}
