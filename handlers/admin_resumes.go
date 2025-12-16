package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterAdminResumeRoutes(r *gin.RouterGroup, db *gorm.DB, auth gin.HandlerFunc, requireAdmin gin.HandlerFunc, rdb *redis.Client) {
	adm := r.Group("/admin")
	adm.Use(auth, requireAdmin)

	adm.GET("/resumes", func(c *gin.Context) {
		page := parseIntDefault(c.Query("page"), 1)
		size := parseIntDefault(c.Query("pageSize"), 20)
		if size > 100 {
			size = 100
		}
		var list []models.Resume
		q := db.Model(&models.Resume{})
		if v := strings.TrimSpace(c.Query("userId")); v != "" {
			q = q.Where("user_id = ?", v)
		}
		if v := strings.TrimSpace(c.Query("title")); v != "" {
			q = q.Where("title LIKE ?", "%"+v+"%")
		}
		if v := strings.TrimSpace(c.Query("templateId")); v != "" {
			q = q.Where("template_id = ?", v)
		}
		var total int64
		q.Count(&total)
		if err := q.Order("updated_at desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		// collect user ids to fetch names
		uidSet := make(map[uint]struct{})
		for _, r := range list {
			if r.UserID != 0 {
				uidSet[r.UserID] = struct{}{}
			}
		}
		uids := make([]uint, 0, len(uidSet))
		for id := range uidSet {
			uids = append(uids, id)
		}
		nameMap := make(map[uint]string, len(uids))
		if len(uids) > 0 {
			var users []models.User
			if err := db.Where("id IN ?", uids).Find(&users).Error; err == nil {
				for _, u := range users {
					if u.Name != "" {
						nameMap[u.ID] = u.Name
					} else {
						nameMap[u.ID] = u.Email
					}
				}
			}
		}
		items := make([]gin.H, 0, len(list))
		for _, r := range list {
			items = append(items, gin.H{
				"id":           r.ExternalID,
				"userId":       r.UserID,
				"userName":     nameMap[r.UserID],
				"title":        r.Title,
				"templateId":   r.TemplateID,
				"themeConfig":  gin.H{"color": r.ThemeColor, "fontFamily": r.ThemeFont, "spacing": r.ThemeSpacing},
				"lastModified": r.LastModified,
				"createdAt":    r.CreatedAt,
				"updatedAt":    r.UpdatedAt,
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": items, "page": page, "pageSize": size, "total": total})
	})

	adm.GET("/resumes/:id", func(c *gin.Context) {
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).Preload("Sections.Items").First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	adm.DELETE("/resumes/:id", func(c *gin.Context) {
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if err := db.Delete(&res).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		writeAudit(c, "resume.delete", "resume", c.Param("id"), "")
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	adm.PATCH("/resumes/:id/visibility", func(c *gin.Context) {
		var body struct {
			IsPublic bool `json:"isPublic"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		var sl models.ShareLink
		if err := db.Where("resume_id = ?", res.ID).First(&sl).Error; err != nil {
			sl = models.ShareLink{ResumeID: res.ID, Slug: uuid.NewString()[:8], IsPublic: body.IsPublic}
			if err := db.Create(&sl).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
				return
			}
		} else {
			sl.IsPublic = body.IsPublic
			if err := db.Save(&sl).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
				return
			}
		}
		if rdb != nil {
			_ = rdb.Del(c, "public:resume:"+sl.Slug).Err()
		}
		writeAudit(c, "resume.visibility", "resume", c.Param("id"), strconv.FormatBool(body.IsPublic))
		c.JSON(http.StatusOK, gin.H{"success": true, "slug": sl.Slug})
	})
}
