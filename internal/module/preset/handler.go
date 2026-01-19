package preset

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

func (h *Handler) GetByID(c *gin.Context) {
	p, err := h.svc.GetByExternalActive(c.Param("id"))
	if err != nil {
		logger.WithCtx(c).Error("preset.get failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, ContentPresetDTO{
		ExternalID:     p.ExternalID,
		Name:           p.Name,
		Language:       p.Language,
		RoleExternalID: p.RoleExternalID,
		Tags:           splitTags(p.Tags),
		DataJSON:       p.DataJSON,
		IsActive:       p.IsActive,
	})
}

