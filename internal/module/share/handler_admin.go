package share

import (
	"net/http"
	"strconv"
	"strings"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/database"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
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
					nameMap[u.ID] = u.Email
				}
			}
		}
	}
	items := make([]gin.H, 0, len(list))
	for _, s := range list {
		uid := s.UserID
		if uid == 0 {
			uid = userIDByResume[s.ResumeID]
		}
		items = append(items, gin.H{
			"id":        s.ID,
			"resumeId":  s.ResumeID,
			"userId":    uid,
			"userName":  nameMap[uid],
			"slug":      s.Slug,
			"url":       "/#/public/" + s.Slug,
			"apiUrl":    "/public/resumes/" + s.Slug,
			"isPublic":  s.IsPublic,
			"createdAt": s.CreatedAt,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "page": page, "pageSize": size, "total": total})
}

func (h *AdminHandler) AdminUpdate(c *gin.Context) {
	var body struct {
		IsPublic bool `json:"isPublic"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var s ShareLink
	if err := database.DB.Where("slug = ?", c.Param("slug")).First(&s).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	s.IsPublic = body.IsPublic
	if err := database.DB.Save(&s).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if cache.RDB != nil {
		_ = cache.RDB.Del(c, common.RedisKeyPublicResume.F(slug)).Err()
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
