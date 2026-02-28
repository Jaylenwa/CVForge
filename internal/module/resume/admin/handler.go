package admin

import (
	"net/http"
	"strconv"
	"strings"

	"cvforge/internal/middleware"
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
	var userID *uint
	if v := strings.TrimSpace(c.Query("userId")); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil && n > 0 {
			x := uint(n)
			userID = &x
		}
	}
	title := strings.TrimSpace(c.Query("title"))
	templateID := strings.TrimSpace(c.Query("templateId"))
	out, err := h.svc.AdminListResumes(page, size, userID, title, templateID)
	if err != nil {
		logger.WithCtx(c).Error("resume.admin_list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) AdminGet(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	out, err := h.svc.AdminGetResume(uint(id))
	if err != nil {
		logger.WithCtx(c).Error("resume.admin_get not found", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *Handler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	actorID, _ := middleware.UID(c)
	actor := AuditActor{ActorID: actorID, IP: c.ClientIP(), UA: c.GetHeader("User-Agent")}
	if err := h.svc.AdminDeleteResume(actor, uint(id)); err != nil {
		logger.WithCtx(c).Error("resume.admin_delete failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminUpdateVisibility(c *gin.Context) {
	var body struct {
		IsPublic bool `json:"isPublic"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	actorID, _ := middleware.UID(c)
	actor := AuditActor{ActorID: actorID, IP: c.ClientIP(), UA: c.GetHeader("User-Agent")}
	slug, err := h.svc.AdminUpdateResumeVisibility(c.Request.Context(), actor, uint(id), body.IsPublic)
	if err != nil {
		logger.WithCtx(c).Error("resume.admin_update_visibility failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "slug": slug})
}

