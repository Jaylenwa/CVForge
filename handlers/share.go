package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterShareRoutes(r *gin.RouterGroup, db *gorm.DB, rdb *redis.Client, auth gin.HandlerFunc) {
	r.POST("/resumes/:id/publish", auth, func(c *gin.Context) {
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
			return
		}
		var sl models.ShareLink
		if err := db.Where("resume_id = ?", res.ID).First(&sl).Error; err != nil {
			sl = models.ShareLink{ResumeID: res.ID, Slug: uuid.NewString()[:8], IsPublic: true}
			_ = db.Create(&sl).Error
		} else {
			sl.IsPublic = true
			_ = db.Save(&sl).Error
		}
		c.JSON(http.StatusOK, gin.H{"slug": sl.Slug, "url": "/public/resumes/" + sl.Slug})
	})

	r.GET("/public/resumes/:slug", func(c *gin.Context) {
		slug := c.Param("slug")
		cacheKey := "public:resume:" + slug
		if val, err := rdb.Get(context.Background(), cacheKey).Result(); err == nil {
			_ = rdb.Incr(context.Background(), "views:"+slug)
			c.Data(http.StatusOK, "application/json", []byte(val))
			return
		}
		var sl models.ShareLink
		if err := db.Where("slug = ? AND is_public = ?", slug, true).First(&sl).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		var res models.Resume
		if err := db.Where("id = ?", sl.ResumeID).Preload("Sections.Items").First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
			return
		}
		payloadObj := gin.H{
			"Title":        res.Title,
			"TemplateID":   res.TemplateID,
			"ThemeColor":   res.ThemeColor,
			"ThemeFont":    res.ThemeFont,
			"ThemeSpacing": res.ThemeSpacing,
			"FullName":     res.FullName,
			"Email":        res.Email,
			"Phone":        res.Phone,
			"Website":      res.Website,
			"AvatarURL":    res.AvatarURL,
			"Sections":     res.Sections,
		}
		payload, _ := json.Marshal(payloadObj)
		_ = rdb.Set(context.Background(), cacheKey, string(payload), 10*time.Minute).Err()
		_ = rdb.Incr(context.Background(), "views:"+slug)
		c.JSON(http.StatusOK, payloadObj)
	})
}
