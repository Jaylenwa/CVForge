package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
