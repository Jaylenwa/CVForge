package template

import (
	"net/http"

	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

func (h *Handler) ListAll(c *gin.Context) {
	payload, err := h.svc.ListAllPayload()
	if err != nil {
		logger.WithCtx(c).Error("template.list_all failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(payload))
}

func (h *Handler) GetByID(c *gin.Context) {
	t, err := h.svc.GetByExternal(c.Param("id"))
	if err != nil {
		logger.WithCtx(c).Error("template.get_by_id failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}
