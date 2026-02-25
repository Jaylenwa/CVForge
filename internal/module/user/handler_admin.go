package user

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/database"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AdminHandler struct {
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func writeAudit(c *gin.Context, action, targetType, targetID, metadata string) {
	actorVal, _ := c.Get("uid")
	ip := c.ClientIP()
	ua := c.GetHeader("User-Agent")
	_ = database.DB.Create(&AuditLog{
		ActorID:    toUint(actorVal),
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   metadata,
		IP:         ip,
		UA:         ua,
	}).Error
}

func toUint(v any) uint {
	switch t := v.(type) {
	case uint:
		return t
	case int:
		if t < 0 {
			return 0
		}
		return uint(t)
	default:
		return 0
	}
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

func (h *AdminHandler) AdminList(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	if size > 100 {
		size = 100
	}
	var list []User
	q := database.DB.Model(&User{})
	emailQ := strings.TrimSpace(c.Query("email"))
	nameQ := strings.TrimSpace(c.Query("name"))
	if emailQ != "" || nameQ != "" {
		if emailQ != "" && nameQ != "" {
			q = q.Where("(email LIKE ? OR name LIKE ?)", "%"+emailQ+"%", "%"+nameQ+"%")
		} else if emailQ != "" {
			q = q.Where("email LIKE ?", "%"+emailQ+"%")
		} else {
			q = q.Where("name LIKE ?", "%"+nameQ+"%")
		}
	}
	if v := strings.TrimSpace(c.Query("role")); v != "" {
		q = q.Where("role = ?", v)
	}
	if v := strings.TrimSpace(c.Query("isActive")); v != "" {
		switch v {
		case "true":
			q = q.Where("is_active = ?", true)
		case "false":
			q = q.Where("is_active = ?", false)
		default:
			logger.WithCtx(c).Error("user.admin_list invalid isActive")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid isActive"})
			return
		}
	}
	var total int64
	q.Count(&total)
	if err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
}

func (h *AdminHandler) AdminGet(c *gin.Context) {
	var u User
	if err := database.DB.First(&u, c.Param("id")).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_get not found", zap.Error(err), zap.String("id", c.Param("id")))
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
}

func (h *AdminHandler) AdminPatch(c *gin.Context) {
	var body struct {
		Name      *string      `json:"name"`
		AvatarURL *string      `json:"avatarUrl"`
		Language  *string      `json:"language"`
		Role      *common.Role `json:"role"`
		IsActive  *bool        `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("user.admin_patch bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var u User
	if err := database.DB.First(&u, c.Param("id")).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_patch not found", zap.Error(err), zap.String("id", c.Param("id")))
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
		r := *body.Role
		if r == common.RoleUser || r == common.RoleAdmin {
			u.Role = r
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
	}
	if body.IsActive != nil {
		u.IsActive = *body.IsActive
	}
	if err := database.DB.Save(&u).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_patch save failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	writeAudit(c, "user.update", "user", strconv.FormatUint(uint64(u.ID), 10), "")
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminResetPassword(c *gin.Context) {
	var body struct {
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.NewPassword == "" {
		logger.WithCtx(c).Error("user.admin_reset_password bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var u User
	if err := database.DB.First(&u, c.Param("id")).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_reset_password not found", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	u.PasswordHash = string(hash)
	if err := database.DB.Save(&u).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_reset_password save failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	writeAudit(c, "user.reset_password", "user", strconv.FormatUint(uint64(u.ID), 10), "")
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminBan(c *gin.Context) {
	if err := database.DB.Model(&User{}).Where("id = ?", c.Param("id")).Update("is_active", false).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_ban failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cache.RDB != nil {
		_ = cache.RDB.Set(context.Background(), common.RedisKeyUserActive.F(c.Param("id")), "0", 10*time.Minute).Err()
	}
	writeAudit(c, "user.ban", "user", c.Param("id"), "")
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminUnban(c *gin.Context) {
	if err := database.DB.Model(&User{}).Where("id = ?", c.Param("id")).Update("is_active", true).Error; err != nil {
		logger.WithCtx(c).Error("user.admin_unban failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cache.RDB != nil {
		_ = cache.RDB.Set(context.Background(), common.RedisKeyUserActive.F(c.Param("id")), "1", 10*time.Minute).Err()
	}
	writeAudit(c, "user.unban", "user", c.Param("id"), "")
	c.JSON(http.StatusOK, gin.H{"success": true})
}
