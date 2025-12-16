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
		{ExternalID: "t1", Name: "Classic Professional", Tags: "Professional,Simple,ATS Friendly", Popularity: 98, IsPremium: false, Category: "General"},
		{ExternalID: "t2", Name: "Modern Dark", Tags: "Creative,Design,Startup", Popularity: 85, IsPremium: true, Category: "Creative"},
		{ExternalID: "t3", Name: "Tech Minimalist", Tags: "Minimalist,Tech,Clean", Popularity: 92, IsPremium: false, Category: "IT"},
		{ExternalID: "t4", Name: "Executive Serif", Tags: "Professional,Management,Senior", Popularity: 70, IsPremium: true, Category: "Finance"},
		{ExternalID: "t5", Name: "Creative Bold", Tags: "Creative,Marketing,Colorful", Popularity: 65, IsPremium: true, Category: "Creative"},
		{ExternalID: "t6", Name: "Elegant Teal", Tags: "Modern,Fresh,Entry Level", Popularity: 88, IsPremium: false, Category: "General"},
		{ExternalID: "t7", Name: "Chinese Blue", Tags: "General,Chinese,ATS Friendly", Popularity: 75, IsPremium: false, Category: "General"},
	}
	for _, m := range mocks {
		_ = db.Create(&m).Error
	}
}
