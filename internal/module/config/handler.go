package config

import (
	"net/http"
	"cvforge/internal/models"
	"cvforge/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	service *Service
}

func NewHandler() *Handler {
	return &Handler{service: NewService()}
}

func (h *Handler) GetPublic(c *gin.Context) {
	c.JSON(http.StatusOK, h.service.GetPublicConfig())
}

func (h *Handler) AdminList(c *gin.Context) {
	configs, err := h.service.GetAll()
	if err != nil {
		logger.WithCtx(c).Error("config.admin_list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, configs)
}

type UpdateRequest struct {
	Configs []models.Config `json:"configs"`
}

func (h *Handler) AdminUpdate(c *gin.Context) {
	var req UpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]string{}
	for _, cfg := range req.Configs {
		updates[cfg.ConfigKey] = cfg.ConfigValue
	}
	// Validate storage configuration strictly
	if v, ok := updates["enabled_storage_s3"]; ok {
		enabled := v == "true" || v == "on" || v == "1"
		if enabled {
			bucket := updates["storage_s3_bucket"]
			region := updates["storage_s3_region"]
			if bucket == "" {
				bucket = h.service.Get("storage_s3_bucket")
			}
			if region == "" {
				region = h.service.Get("storage_s3_region")
			}
			if bucket == "" || region == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "S3 enabled requires bucket and region"})
				return
			}
			ak := updates["storage_s3_access_key"]
			sk := updates["storage_s3_secret_key"]
			if ak == "" {
				ak = h.service.Get("storage_s3_access_key")
			}
			if sk == "" {
				sk = h.service.Get("storage_s3_secret_key")
			}
			if (ak != "" && sk == "") || (ak == "" && sk != "") {
				c.JSON(http.StatusBadRequest, gin.H{"error": "S3 credentials must provide both access_key and secret_key"})
				return
			}
		}
	}

	for _, cfg := range req.Configs {
		if err := h.service.Set(cfg.ConfigKey, cfg.ConfigValue, cfg.Description, cfg.Type); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
