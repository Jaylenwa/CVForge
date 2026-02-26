package share

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	"cvforge/internal/infra/database"
	"cvforge/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
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
	var list []ShareLink
	q := database.DB.Model(&ShareLink{})
	if v := strings.TrimSpace(c.Query("slug")); v != "" {
		q = q.Where("slug LIKE ?", "%"+v+"%")
	}
	if v := strings.TrimSpace(c.Query("isPublic")); v != "" {
		if v == "true" {
			q = q.Where("is_public = ?", true)
		} else if v == "false" {
			q = q.Where("is_public = ?", false)
		}
	}
	var total int64
	q.Count(&total)
	if err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		logger.WithCtx(c).Error("share.admin_list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Map resume_id -> user_id
	resumeIDs := make([]uint, 0, len(list))
	for _, s := range list {
		if s.ResumeID != 0 {
			resumeIDs = append(resumeIDs, s.ResumeID)
		}
	}
	userIDByResume := make(map[uint]uint)
	if len(resumeIDs) > 0 {
		var rs []Resume
		if err := database.DB.Where("id IN ?", resumeIDs).Find(&rs).Error; err == nil {
			for _, r := range rs {
				userIDByResume[r.ID] = r.UserID
			}
		}
	}
	// Collect user IDs
	uidSet := make(map[uint]struct{})
	for _, s := range list {
		uid := s.UserID
		if uid == 0 {
			uid = userIDByResume[s.ResumeID]
		}
		if uid != 0 {
			uidSet[uid] = struct{}{}
		}
	}
	uids := make([]uint, 0, len(uidSet))
	for id := range uidSet {
		uids = append(uids, id)
	}
	nameMap := make(map[uint]string, len(uids))
	if len(uids) > 0 {
		var users []User
		if err := database.DB.Where("id IN ?", uids).Find(&users).Error; err == nil {
			for _, u := range users {
				if u.Name != "" {
					nameMap[u.ID] = u.Name
				} else {
					if u.Email != nil {
						nameMap[u.ID] = *u.Email
					}
				}
			}
		}
	}
	items := make([]gin.H, 0, len(list))
	viewsBySlug := map[string]uint64{}
	lastAccessBySlug := map[string]time.Time{}
	if cache.RDB != nil && len(list) > 0 {
		viewKeys := make([]string, 0, len(list))
		lastKeys := make([]string, 0, len(list))
		for _, s := range list {
			viewKeys = append(viewKeys, common.RedisKeyViews.F(s.Slug))
			lastKeys = append(lastKeys, common.RedisKeyShareLastAccess.F(s.Slug))
		}
		if vals, err := cache.RDB.MGet(c, viewKeys...).Result(); err == nil {
			for i, v := range vals {
				if v == nil || i >= len(list) {
					continue
				}
				switch t := v.(type) {
				case string:
					if n, err := strconv.ParseUint(t, 10, 64); err == nil {
						viewsBySlug[list[i].Slug] = n
					}
				}
			}
		}
		if vals, err := cache.RDB.MGet(c, lastKeys...).Result(); err == nil {
			for i, v := range vals {
				if v == nil || i >= len(list) {
					continue
				}
				switch t := v.(type) {
				case string:
					if ms, err := strconv.ParseInt(t, 10, 64); err == nil && ms > 0 {
						lastAccessBySlug[list[i].Slug] = time.UnixMilli(ms)
					}
				}
			}
		}
	}
	for _, s := range list {
		uid := s.UserID
		if uid == 0 {
			uid = userIDByResume[s.ResumeID]
		}
		views := s.Views
		if v, ok := viewsBySlug[s.Slug]; ok && v > 0 {
			views = v
		}
		var lastAccessAt any = s.LastAccessAt
		if t, ok := lastAccessBySlug[s.Slug]; ok {
			tt := t
			lastAccessAt = &tt
		}
		items = append(items, gin.H{
			"id":           s.ID,
			"resumeId":     s.ResumeID,
			"userId":       uid,
			"userName":     nameMap[uid],
			"slug":         s.Slug,
			"url":          "/#/public/" + s.Slug,
			"apiUrl":       "/public/resumes/" + s.Slug,
			"isPublic":     s.IsPublic,
			"hasPassword":  s.Password != "",
			"expiresAt":    s.ExpiresAt,
			"views":        views,
			"lastAccessAt": lastAccessAt,
			"createdAt":    s.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "page": page, "pageSize": size, "total": total})
}

func (h *AdminHandler) AdminUpdate(c *gin.Context) {
	var body struct {
		IsPublic bool `json:"isPublic"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		logger.WithCtx(c).Error("share.admin_update bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var s ShareLink
	if err := database.DB.Where("slug = ?", c.Param("slug")).First(&s).Error; err != nil {
		logger.WithCtx(c).Error("share.admin_update not found", zap.Error(err), zap.String("slug", c.Param("slug")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	s.IsPublic = body.IsPublic
	if err := database.DB.Save(&s).Error; err != nil {
		logger.WithCtx(c).Error("share.admin_update save failed", zap.Error(err), zap.String("slug", c.Param("slug")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cache.RDB != nil {
		_ = cache.RDB.Del(c, common.RedisKeyPublicResume.F(s.Slug)).Err()
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDelete(c *gin.Context) {
	slug := c.Param("slug")
	if err := database.DB.Where("slug = ?", slug).Delete(&ShareLink{}).Error; err != nil {
		logger.WithCtx(c).Error("share.admin_delete failed", zap.Error(err), zap.String("slug", slug))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if cache.RDB != nil {
		_ = cache.RDB.Del(c, common.RedisKeyPublicResume.F(slug)).Err()
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
