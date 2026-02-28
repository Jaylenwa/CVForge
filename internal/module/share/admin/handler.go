package admin

import (
	"net/http"
	"strconv"
	"strings"

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
	if size > 100 {
		size = 100
	}
	slug := strings.TrimSpace(c.Query("slug"))
	var isPublic *bool
	if v := strings.TrimSpace(c.Query("isPublic")); v != "" {
		if v == "true" {
			b := true
			isPublic = &b
		} else if v == "false" {
			b := false
			isPublic = &b
		}
	}
	list, total, err := h.svc.AdminListShareLinks(c.Request.Context(), page, size, slug, isPublic)
	if err != nil {
		logger.WithCtx(c).Error("share.admin_list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	items := make([]gin.H, 0, len(list))
	for _, it := range list {
		items = append(items, gin.H{
			"id":           it.ID,
			"resumeId":     it.ResumeID,
			"userId":       it.UserID,
			"userName":     it.UserName,
			"slug":         it.Slug,
			"url":          "/#/public/" + it.Slug,
			"apiUrl":       "/public/resumes/" + it.Slug,
			"isPublic":     it.IsPublic,
			"hasPassword":  it.HasPassword,
			"expiresAt":    it.ExpiresAt,
			"views":        it.Views,
			"lastAccessAt": it.LastAccessAt,
			"createdAt":    it.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "page": page, "pageSize": size, "total": total})
}

func (h *Handler) AdminUpdate(c *gin.Context) {
	var body struct {
		IsPublic bool `json:"isPublic"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("share.admin_update bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := h.svc.AdminUpdateShareLink(c.Request.Context(), c.Param("slug"), body.IsPublic); err != nil {
		logger.WithCtx(c).Error("share.admin_update failed", zap.Error(err), zap.String("slug", c.Param("slug")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) AdminDelete(c *gin.Context) {
	slug := c.Param("slug")
	if err := h.svc.AdminDeleteShareLink(c.Request.Context(), slug); err != nil {
		logger.WithCtx(c).Error("share.admin_delete failed", zap.Error(err), zap.String("slug", slug))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

