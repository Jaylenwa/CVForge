package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterTemplateRoutes(r *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	r.GET("/templates", func(c *gin.Context) {
		if rdb != nil {
			if val, err := rdb.Get(context.Background(), "templates:list:all").Result(); err == nil {
				c.Data(http.StatusOK, "application/json", []byte(val))
				return
			}
		}
		var count int64
		db.Model(&models.Template{}).Count(&count)
		if count == 0 {
			seedTemplates(db)
		}
		var list []models.Template
		if err := db.Order("popularity desc").Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		payload, _ := json.Marshal(gin.H{"items": list})
		if rdb != nil {
			_ = rdb.Set(context.Background(), "templates:list:all", string(payload), time.Hour).Err()
		}
		c.Data(http.StatusOK, "application/json", payload)
	})
	r.GET("/templates/:id", func(c *gin.Context) {
		var t models.Template
		if err := db.Where("external_id = ?", c.Param("id")).First(&t).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, t)
	})
}

func seedTemplates(db *gorm.DB) {
	mocks := []models.Template{
		{ExternalID: "t1", Name: "Classic Professional", Thumbnail: "https://images.unsplash.com/photo-1586281380349-632531db7ed4?w=500&q=80", Tags: "Professional,Simple,ATS Friendly", Popularity: 98, IsPremium: false, Category: "General", Level: "Senior"},
		{ExternalID: "t2", Name: "Modern Dark", Thumbnail: "https://images.unsplash.com/photo-1616628188859-7a11abb6fcc9?w=500&q=80", Tags: "Creative,Design,Startup", Popularity: 85, IsPremium: true, Category: "Creative", Level: "Junior"},
		{ExternalID: "t3", Name: "Tech Minimalist", Thumbnail: "https://images.unsplash.com/photo-1512486130939-2c4f79935e4f?w=500&q=80", Tags: "Minimalist,Tech,Clean", Popularity: 92, IsPremium: false, Category: "IT", Level: "Junior"},
		{ExternalID: "t4", Name: "Executive Serif", Thumbnail: "https://images.unsplash.com/photo-1507679799987-c73779587ccf?w=500&q=80", Tags: "Professional,Management,Senior", Popularity: 70, IsPremium: true, Category: "Finance", Level: "Executive"},
		{ExternalID: "t5", Name: "Creative Bold", Thumbnail: "https://images.unsplash.com/photo-1513542789411-b6a5d4f31634?w=500&q=80", Tags: "Creative,Marketing,Colorful", Popularity: 65, IsPremium: true, Category: "Creative", Level: "Senior"},
		{ExternalID: "t6", Name: "Elegant Teal", Thumbnail: "https://images.unsplash.com/photo-1515378791036-0648a3ef77b2?w=500&q=80", Tags: "Modern,Fresh,Entry Level", Popularity: 88, IsPremium: false, Category: "General", Level: "Intern"},
	}
	for _, m := range mocks {
		_ = db.Create(&m).Error
	}
}
