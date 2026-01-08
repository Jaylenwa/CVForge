package ai

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

type polishReq struct {
	Text string `json:"text"`
	Tone string `json:"tone"`
}
type summaryReq struct {
	Job    string `json:"job"`
	Skills string `json:"skills"`
}

func (h *Handler) Polish(c *gin.Context) {
	var req polishReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	out := h.svc.Polish(req.Text, req.Tone)
	c.JSON(http.StatusOK, gin.H{"text": out})
}

func (h *Handler) Summary(c *gin.Context) {
	var req summaryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	text := h.svc.Summary(req.Job, req.Skills)
	c.JSON(http.StatusOK, gin.H{"text": text})
}
