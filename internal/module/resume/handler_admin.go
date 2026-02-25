package resume

import (
	"net/http"
	"strconv"
	"strings"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/database"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AdminHandler struct {
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func (h *AdminHandler) AdminList(c *gin.Context) {
	page := parseIntDefault(c.Query("page"), 1)
	size := parseIntDefault(c.Query("pageSize"), 20)
	if size > 100 {
		size = 100
	}
	var list []Resume
	q := database.DB.Model(&Resume{})
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
	if err := q.Preload("Personal").Preload("Theme").Order("updated_at desc").Offset((page - 1) * size).Limit(size).Find(&list).Error; err != nil {
		logger.WithCtx(c).Error("resume.admin_list failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
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
	for _, r := range list {
		items = append(items, gin.H{
			"resume":    ToDTO(r),
			"userId":    r.UserID,
			"userName":  nameMap[r.UserID],
			"createdAt": r.CreatedAt,
			"updatedAt": r.UpdatedAt,
		})
	}
	totalPages := (int(total) + size - 1) / size
	hasNext := page*size < int(total)
	c.JSON(http.StatusOK, gin.H{
		"items":      items,
		"page":       page,
		"pageSize":   size,
		"total":      total,
		"totalPages": totalPages,
		"hasNext":    hasNext,
	})
}

func (h *AdminHandler) AdminGet(c *gin.Context) {
	var res Resume
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := database.DB.Where("id = ?", uint(id)).Preload("Personal").Preload("Theme").Preload("Sections.Items").First(&res).Error; err != nil {
		logger.WithCtx(c).Error("resume.admin_get not found", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, ToDTO(res))
}

func (h *AdminHandler) AdminDelete(c *gin.Context) {
	var res Resume
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := database.DB.Where("id = ?", uint(id)).First(&res).Error; err != nil {
		logger.WithCtx(c).Error("resume.admin_delete not found", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := database.DB.Delete(&res).Error; err != nil {
		logger.WithCtx(c).Error("resume.admin_delete failed", zap.Error(err), zap.String("id", c.Param("id")))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	writeAudit(c, "resume.delete", "resume", c.Param("id"), "")
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminUpdateVisibility(c *gin.Context) {
	var body struct {
		IsPublic bool `json:"isPublic"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	var res Resume
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	if err := database.DB.Where("id = ?", uint(id)).First(&res).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var sl ShareLink
	if err := database.DB.Where("resume_id = ?", res.ID).First(&sl).Error; err != nil {
		sl = ShareLink{ResumeID: res.ID, Slug: uuid.NewString()[:8], IsPublic: body.IsPublic}
		if err := database.DB.Create(&sl).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		sl.IsPublic = body.IsPublic
		if err := database.DB.Save(&sl).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	if cache.RDB != nil {
		_ = cache.RDB.Del(c, common.RedisKeyPublicResume.F(sl.Slug)).Err()
	}
	writeAudit(c, "resume.visibility", "resume", c.Param("id"), strconv.FormatBool(body.IsPublic))
	c.JSON(http.StatusOK, gin.H{"success": true, "slug": sl.Slug})
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
