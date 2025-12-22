package config

import (
	"net/http"
	"openresume/internal/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetPublic(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetPublicConfig())
}

func (h *Handler) AdminList(c *gin.Context) {
	configs, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch configs"})
		return
	}
	c.JSON(http.StatusOK, configs)
}

type UpdateRequest struct {
	Configs []models.SystemConfig `json:"configs"`
}

func (h *Handler) AdminUpdate(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, cfg := range req.Configs {
		if err := h.service.Set(cfg.Key, cfg.Value, cfg.Description, cfg.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update config: " + cfg.Key})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
