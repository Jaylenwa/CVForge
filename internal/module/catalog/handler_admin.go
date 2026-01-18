package catalog

import (
	"net/http"

	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler struct {
	svc *Service
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{svc: NewService()}
}

func (h *AdminHandler) AdminImportSeed(c *gin.Context) {
	seed, err := DefaultSeed()
	if err != nil {
		logger.WithCtx(c).Error("catalog.seed.build failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "seed error"})
		return
	}
	counts, err := h.svc.ImportSeed(seed)
	if err != nil {
		logger.WithCtx(c).Error("catalog.seed.import failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "counts": counts})
}

