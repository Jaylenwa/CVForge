package handlers

import (
	"net/http"

	"openresume/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUserRoutesReal(r *gin.RouterGroup, db *gorm.DB, auth gin.HandlerFunc) {
	r.GET("/users/me", auth, func(c *gin.Context) {
		uidVal, _ := c.Get("uid")
		var u models.User
		if err := db.First(&u, uidVal).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"email":     u.Email,
			"name":      u.Name,
			"avatarUrl": u.AvatarURL,
			"language":  u.Language,
			"role":      u.Role,
		})
	})

	r.PUT("/users/profile", auth, func(c *gin.Context) {
		uidVal, _ := c.Get("uid")
		var body struct{ Name, AvatarURL, Language string }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var u models.User
		if err := db.First(&u, uidVal).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		u.Name = ifNotEmpty(body.Name, u.Name)
		u.AvatarURL = ifNotEmpty(body.AvatarURL, u.AvatarURL)
		u.Language = ifNotEmpty(body.Language, u.Language)
		if err := db.Save(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.PUT("/users/password", auth, func(c *gin.Context) {
		uidVal, _ := c.Get("uid")
		var body struct{ CurrentPassword, NewPassword string }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var u models.User
		if err := db.First(&u, uidVal).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		if body.NewPassword == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid new password"})
			return
		}
		// 仅验证存在旧密码哈希，不强制旧密码为必填（演示）
		if u.PasswordHash != "" {
			if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(body.CurrentPassword)); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid current password"})
				return
			}
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
		u.PasswordHash = string(hash)
		if err := db.Save(&u).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}

func ifNotEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
