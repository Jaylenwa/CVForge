package stats

import (
	"net/http"
	"strconv"
	"time"

	"openresume/internal/common"
	"openresume/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"context"
)

type Handler struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewHandler(db *gorm.DB, rdb *redis.Client) *Handler {
	return &Handler{db: db, rdb: rdb}
}

type StatsResponse struct {
	Totals struct {
		Users     int64 `json:"users"`
		Resumes   int64 `json:"resumes"`
		Templates int64 `json:"templates"`
		VisitorsToday int64 `json:"visitorsToday"`
	} `json:"totals"`
	Trend struct {
		Dates     []string `json:"dates"`
		Users     []int64  `json:"users"`
		Resumes   []int64  `json:"resumes"`
		Templates []int64  `json:"templates"`
		Visitors  []int64  `json:"visitors"`
	} `json:"trend"`
	GeneratedAt int64 `json:"generatedAt"`
}

func (h *Handler) AdminStats(c *gin.Context) {
	days := 14
	if v := c.Query("days"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 60 {
			days = n
		}
	}

	var userTotal int64
	var resumeTotal int64
	var templateTotal int64
	_ = h.db.Model(&models.User{}).Count(&userTotal).Error
	_ = h.db.Model(&models.Resume{}).Count(&resumeTotal).Error
	_ = h.db.Model(&models.Template{}).Count(&templateTotal).Error

	now := time.Now()
	// Truncate to local midnight
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Add(24 * time.Hour)

	dates := make([]string, 0, days)
	users := make([]int64, 0, days)
	resumes := make([]int64, 0, days)
	templates := make([]int64, 0, days)
	visitors := make([]int64, 0, days)

	for i := days - 1; i >= 0; i-- {
		dayStart := end.Add(-time.Duration(i+1) * 24 * time.Hour)
		dayEnd := end.Add(-time.Duration(i) * 24 * time.Hour)

		dates = append(dates, dayStart.Format("2006-01-02"))

		var uc int64
		var rc int64
		var tc int64
		_ = h.db.Model(&models.User{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&uc).Error
		_ = h.db.Model(&models.Resume{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&rc).Error
		_ = h.db.Model(&models.Template{}).Where("created_at >= ? AND created_at < ?", dayStart, dayEnd).Count(&tc).Error
		var vc int64
		if h.rdb != nil {
			key := common.RedisKeyUVDay.F(dayStart.Format("2006-01-02"))
			if n, err := h.rdb.SCard(context.Background(), key).Result(); err == nil {
				vc = n
			}
		}

		users = append(users, uc)
		resumes = append(resumes, rc)
		templates = append(templates, tc)
		visitors = append(visitors, vc)
	}

	out := StatsResponse{
		GeneratedAt: now.Unix(),
	}
	out.Totals.Users = userTotal
	out.Totals.Resumes = resumeTotal
	out.Totals.Templates = templateTotal
	out.Trend.Dates = dates
	out.Trend.Users = users
	out.Trend.Resumes = resumes
	out.Trend.Templates = templates
	out.Trend.Visitors = visitors
	if h.rdb != nil {
		todayKey := common.RedisKeyUVDay.F(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format("2006-01-02"))
		if n, err := h.rdb.SCard(context.Background(), todayKey).Result(); err == nil {
			out.Totals.VisitorsToday = n
		}
	}

	c.JSON(http.StatusOK, out)
}
