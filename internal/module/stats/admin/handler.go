package admin

import (
	"net/http"
	"strconv"

	statsmod "cvforge/internal/module/stats"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *statsmod.Service
}

func NewHandler() *Handler {
	return &Handler{svc: statsmod.NewService()}
}

func (h *Handler) AdminStats(c *gin.Context) {
	days := 14
	if v := c.Query("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 60 {
			days = n
		}
	}
	out, err := h.svc.AdminStats(days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

