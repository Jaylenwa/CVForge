package admin

import (
	"net/http"
	"strconv"
	"strings"

	"cvforge/internal/common"
	"cvforge/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
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

func (h *Handler) AdminList(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	if size > 100 {
		size = 100
	}
	emailQ := strings.TrimSpace(c.Query("email"))
	nameQ := strings.TrimSpace(c.Query("name"))
	role := strings.TrimSpace(c.Query("role"))
	var isActive *bool
	if v := strings.TrimSpace(c.Query("isActive")); v != "" {
		switch v {
		case "true":
			b := true
			isActive = &b
		case "false":
			b := false
			isActive = &b
		default:
			logger.WithCtx(c).Error("user.admin_list invalid isActive")
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid isActive"})
			return
		}
	}
	items, total, err := h.svc.AdminListUsers(page, size, emailQ, nameQ, role, isActive)
	if err != nil {
		logger.WithCtx(c).Error("user.admin_list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "page": page, "pageSize": size, "total": total})
}

func (h *Handler) AdminGet(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	out, err := h.svc.AdminGetUser(uint(id))
	if err != nil {
		logger.WithCtx(c).Error("user.admin_get not found", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) AdminPatch(c *gin.Context) {
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
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	actorVal, _ := c.Get("uid")
	actor := AuditActor{ActorID: toUint(actorVal), IP: c.ClientIP(), UA: c.GetHeader("User-Agent")}
	if err := h.svc.AdminPatchUser(actor, uint(id), AdminPatchInput{
		Name:      body.Name,
		AvatarURL: body.AvatarURL,
		Language:  body.Language,
		Role:      body.Role,
		IsActive:  body.IsActive,
	}); err != nil {
		logger.WithCtx(c).Error("user.admin_patch failed", zap.Error(err), zap.String("id", c.Param("id")))
		if err.Error() == "invalid role" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminResetPassword(c *gin.Context) {
	var body struct {
		NewPassword string `json:"newPassword"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.NewPassword == "" {
		logger.WithCtx(c).Error("user.admin_reset_password bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	actorVal, _ := c.Get("uid")
	actor := AuditActor{ActorID: toUint(actorVal), IP: c.ClientIP(), UA: c.GetHeader("User-Agent")}
	if err := h.svc.AdminResetPassword(actor, uint(id), body.NewPassword); err != nil {
		logger.WithCtx(c).Error("user.admin_reset_password failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminBan(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	actorVal, _ := c.Get("uid")
	actor := AuditActor{ActorID: toUint(actorVal), IP: c.ClientIP(), UA: c.GetHeader("User-Agent")}
	if err := h.svc.AdminSetActive(actor, uint(id), false); err != nil {
		logger.WithCtx(c).Error("user.admin_ban failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminUnban(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	actorVal, _ := c.Get("uid")
	actor := AuditActor{ActorID: toUint(actorVal), IP: c.ClientIP(), UA: c.GetHeader("User-Agent")}
	if err := h.svc.AdminSetActive(actor, uint(id), true); err != nil {
		logger.WithCtx(c).Error("user.admin_unban failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
