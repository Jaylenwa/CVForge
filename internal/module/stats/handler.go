package stats

import (
	"net/http"
	"strconv"
	"time"

	"openresume/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

type StatsResponse struct {
	Totals struct {
		Users     int64 `json:"users"`
		Resumes   int64 `json:"resumes"`
		Templates int64 `json:"templates"`
	} `json:"totals"`
	Trend struct {
		Dates     []string `json:"dates"`
		Users     []int64  `json:"users"`
		Resumes   []int64  `json:"resumes"`
		Templates []int64  `json:"templates"`
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

		users = append(users, uc)
		resumes = append(resumes, rc)
		templates = append(templates, tc)
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

	c.JSON(http.StatusOK, out)
}
