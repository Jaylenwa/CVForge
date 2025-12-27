package pdf

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

func (h *Handler) GeneratePDF(c *gin.Context) {
	pdf, code, err := h.svc.GeneratePDF(c, c.Param("id"))
	if err != nil {
		switch code {
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case 503:
			c.JSON(http.StatusNotImplemented, gin.H{"error": "pdf service unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "pdf service unavailable"})
		}
		return
	}
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=resume.pdf")
	c.Writer.Write(pdf)
}

func (h *Handler) GenerateImage(c *gin.Context) {
	img, code, err := h.svc.GenerateImage(c, c.Param("id"))
	if err != nil {
		switch code {
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case 503:
			c.JSON(http.StatusNotImplemented, gin.H{"error": "image service unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "image service unavailable"})
		}
		return
	}
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename=resume.png")
	c.Writer.Write(img)
}
