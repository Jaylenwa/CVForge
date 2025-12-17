package template

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AdminHandler struct {
	svc *Service
}

func NewAdminHandler(db *gorm.DB, rdb *redis.Client) *AdminHandler {
	return &AdminHandler{svc: NewService(db, rdb)}
}

func (h *AdminHandler) AdminCreate(c *gin.Context) {
	var body struct {
		ExternalID string `json:"externalId"`
		Name       string `json:"name"`
		Tags       string `json:"tags"`
		Popularity *int   `json:"popularity"`
		IsPremium  *bool  `json:"isPremium"`
		Category   string `json:"category"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.ExternalID == "" || body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	mt := Template{
		ExternalID: body.ExternalID,
		Name:       body.Name,
		Tags:       body.Tags,
		Category:   body.Category,
	}
	if body.Popularity != nil {
		mt.Popularity = *body.Popularity
	}
	if body.IsPremium != nil {
		mt.IsPremium = *body.IsPremium
	}
	if err := h.svc.Create(mt); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminPatch(c *gin.Context) {
	var body struct {
		Name       *string `json:"name"`
		Tags       *string `json:"tags"`
		Popularity *int    `json:"popularity"`
		IsPremium  *bool   `json:"isPremium"`
		Category   *string `json:"category"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	err := h.svc.Update(c.Param("id"), func(t *Template) {
		if body.Name != nil {
			t.Name = strings.TrimSpace(*body.Name)
		}
		if body.Tags != nil {
			t.Tags = strings.TrimSpace(*body.Tags)
		}
		if body.Popularity != nil {
			t.Popularity = *body.Popularity
		}
		if body.IsPremium != nil {
			t.IsPremium = *body.IsPremium
		}
		if body.Category != nil {
			t.Category = strings.TrimSpace(*body.Category)
		}
	})
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDelete(c *gin.Context) {
	if err := h.svc.Delete(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
