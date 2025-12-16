package handlers

import (
	"net/http"
	"strings"

	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterAdminTemplateRoutes(r *gin.RouterGroup, db *gorm.DB, auth gin.HandlerFunc, requireAdmin gin.HandlerFunc, rdb *redis.Client) {
	adm := r.Group("/admin")
	adm.Use(auth, requireAdmin)

	adm.POST("/templates", func(c *gin.Context) {
		var body struct {
			ExternalID string  `json:"externalId"`
			Name       string  `json:"name"`
			Tags       string  `json:"tags"`
			Popularity *int    `json:"popularity"`
			IsPremium  *bool   `json:"isPremium"`
			Category   string  `json:"category"`
			Level      string  `json:"level"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.ExternalID == "" || body.Name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		t := models.Template{
			ExternalID: body.ExternalID,
			Name:       body.Name,
			Tags:       body.Tags,
			Category:   body.Category,
			Level:      body.Level,
		}
		if body.Popularity != nil {
			t.Popularity = *body.Popularity
		}
		if body.IsPremium != nil {
			t.IsPremium = *body.IsPremium
		}
		if err := db.Create(&t).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
			return
		}
		if rdb != nil {
			_ = rdb.Del(c, "templates:list:all").Err()
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	adm.PATCH("/templates/:id", func(c *gin.Context) {
		var body struct {
			Name       *string `json:"name"`
			Tags       *string `json:"tags"`
			Popularity *int    `json:"popularity"`
			IsPremium  *bool   `json:"isPremium"`
			Category   *string `json:"category"`
			Level      *string `json:"level"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var t models.Template
		if err := db.Where("external_id = ?", c.Param("id")).First(&t).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
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
		if body.Level != nil {
			t.Level = strings.TrimSpace(*body.Level)
		}
		if err := db.Save(&t).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		if rdb != nil {
			_ = rdb.Del(c, "templates:list:all").Err()
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	adm.DELETE("/templates/:id", func(c *gin.Context) {
		if err := db.Where("external_id = ?", c.Param("id")).Delete(&models.Template{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		if rdb != nil {
			_ = rdb.Del(c, "templates:list:all").Err()
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}
