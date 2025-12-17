package share

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AdminHandler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewAdminHandler(db *gorm.DB, rdb *redis.Client) *AdminHandler {
	return &AdminHandler{db: db, rdb: rdb}
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
	q := h.db.Model(&ShareLink{})
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
	items := make([]gin.H, 0, len(list))
	for _, s := range list {
		items = append(items, gin.H{
			"id":        s.ID,
			"resumeId":  s.ResumeID,
			"slug":      s.Slug,
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
	if err := h.db.Where("slug = ?", c.Param("slug")).First(&s).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	s.IsPublic = body.IsPublic
	if err := h.db.Save(&s).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if h.rdb != nil {
		_ = h.rdb.Del(c, "public:resume:"+s.Slug).Err()
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *AdminHandler) AdminDelete(c *gin.Context) {
	slug := c.Param("slug")
	if err := h.db.Where("slug = ?", slug).Delete(&ShareLink{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	if h.rdb != nil {
		_ = h.rdb.Del(c, "public:resume:"+slug).Err()
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
