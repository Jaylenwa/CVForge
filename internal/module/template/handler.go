package template

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

func (h *Handler) ListAll(c *gin.Context) {
	payload, err := h.svc.ListAllPayload()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.Data(http.StatusOK, "application/json", []byte(payload))
}

func (h *Handler) GetByID(c *gin.Context) {
	t, err := h.svc.GetByExternal(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, t)
}
