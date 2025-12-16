package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"openresume/db"
	"openresume/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterAdminUserRoutes(r *gin.RouterGroup, db *gorm.DB, auth gin.HandlerFunc, requireAdmin gin.HandlerFunc) {
	adm := r.Group("/admin")
	adm.Use(auth, requireAdmin)

	adm.GET("/users", func(c *gin.Context) {
		page := parseIntDefault(c.Query("page"), 1)
		size := parseIntDefault(c.Query("pageSize"), 20)
		if size > 100 {
			size = 100
		}
		var list []models.User
		q := db.Model(&models.User{})
		if v := strings.TrimSpace(c.Query("email")); v != "" {
			q = q.Where("email LIKE ?", "%"+v+"%")
		}
		if v := strings.TrimSpace(c.Query("name")); v != "" {
			q = q.Where("name LIKE ?", "%"+v+"%")
		}
		if v := strings.TrimSpace(c.Query("role")); v != "" {
			q = q.Where("role = ?", v)
		}
		if v := strings.TrimSpace(c.Query("isActive")); v != "" {
			if v == "true" {
				q = q.Where("is_active = ?", true)
			} else if v == "false" {
				q = q.Where("is_active = ?", false)
			}
		}
		var total int64
		q.Count(&total)
		if err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		items := make([]gin.H, 0, len(list))
		for _, u := range list {
			items = append(items, gin.H{
				"id":          u.ID,
				"email":       u.Email,
				"name":        u.Name,
				"avatarUrl":   u.AvatarURL,
				"language":    u.Language,
				"role":        u.Role,
				"isActive":    u.IsActive,
				"lastLoginAt": u.LastLoginAt,
				"createdAt":   u.CreatedAt,
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": items, "page": page, "pageSize": size, "total": total})
	})

	adm.GET("/users/:id", func(c *gin.Context) {
		var u models.User
		if err := db.First(&u, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"id":          u.ID,
			"email":       u.Email,
			"name":        u.Name,
			"avatarUrl":   u.AvatarURL,
			"language":    u.Language,
			"role":        u.Role,
			"isActive":    u.IsActive,
			"lastLoginAt": u.LastLoginAt,
			"createdAt":   u.CreatedAt,
			"updatedAt":   u.UpdatedAt,
		})
	})

	adm.PATCH("/users/:id", func(c *gin.Context) {
		var body struct {
			Name      *string `json:"name"`
			AvatarURL *string `json:"avatarUrl"`
			Language  *string `json:"language"`
			Role      *string `json:"role"`
			IsActive  *bool   `json:"isActive"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var u models.User
		if err := db.First(&u, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if body.Name != nil {
			u.Name = *body.Name
		}
		if body.AvatarURL != nil {
			u.AvatarURL = *body.AvatarURL
		}
		if body.Language != nil {
			u.Language = *body.Language
		}
		if body.Role != nil {
			r := strings.ToLower(*body.Role)
			if r == "user" || r == "moderator" || r == "admin" {
				u.Role = r
			} else {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
				return
			}
		}
		if body.IsActive != nil {
			u.IsActive = *body.IsActive
		}
		if err := db.Save(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		writeAudit(c, "user.update", "user", strconv.FormatUint(uint64(u.ID), 10), "")
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	adm.POST("/users/:id/reset-password", func(c *gin.Context) {
		var body struct {
			NewPassword string `json:"newPassword"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.NewPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var u models.User
		if err := db.First(&u, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		u.PasswordHash = string(hash)
		if err := db.Save(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		writeAudit(c, "user.reset_password", "user", strconv.FormatUint(uint64(u.ID), 10), "")
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	adm.POST("/users/:id/ban", func(c *gin.Context) {
		if err := db.Model(&models.User{}).Where("id = ?", c.Param("id")).Update("is_active", false).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		writeAudit(c, "user.ban", "user", c.Param("id"), "")
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	adm.POST("/users/:id/unban", func(c *gin.Context) {
		if err := db.Model(&models.User{}).Where("id = ?", c.Param("id")).Update("is_active", true).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		writeAudit(c, "user.unban", "user", c.Param("id"), "")
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}

func writeAudit(c *gin.Context, action, targetType, targetID, metadata string) {
	actorVal, _ := c.Get("uid")
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")
	_ = db.Gorm().Create(&models.AuditLog{
		ActorID:    toUint(actorVal),
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   metadata,
		IP:         ip,
		UA:         ua,
	}).Error
}

func parseIntDefault(s string, d int) int {
	if s == "" {
		return d
	}
	n, err := strconv.Atoi(s)
	if err != nil || n <= 0 {
		return d
	}
	return n
}

func toUint(v interface{}) uint {
	switch t := v.(type) {
	case uint:
		return t
	case int:
		return uint(t)
	case float64:
		return uint(t)
	default:
		return 0
	}
}
