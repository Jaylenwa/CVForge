package preset

import (
	"net/http"
	"strconv"
	"strings"

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
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	lang := c.Query("language")
	p, err := h.svc.GetByIDActive(uint(id), lang)
	if err != nil {
		logger.WithCtx(c).Error("preset.get failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, ContentPresetDTO{
		ID:       p.ID,
		Name:     p.Name,
		Language: p.Language,
		RoleID:   p.RoleID,
		DataJSON: p.DataJSON,
		IsActive: p.IsActive,
	})
}
