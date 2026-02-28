package stats

import (
	"context"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

type StatsResponse struct {
	Totals struct {
		Users         int64 `json:"users"`
		Resumes       int64 `json:"resumes"`
		Templates     int64 `json:"templates"`
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

func (s *Service) AdminStats(days int) (StatsResponse, error) {
	if days <= 0 {
		days = 14
	}
	if days > 60 {
		days = 60
	}

	userTotal, _ := s.repo.CountUsers()
	resumeTotal, _ := s.repo.CountResumes()
	templateTotal, _ := s.repo.CountTemplates()

	now := time.Now()
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

		uc, _ := s.repo.CountUsersCreatedBetween(dayStart, dayEnd)
		rc, _ := s.repo.CountResumesCreatedBetween(dayStart, dayEnd)
		tc, _ := s.repo.CountTemplatesCreatedBetween(dayStart, dayEnd)
		var vc int64
		if cache.RDB != nil {
			key := common.RedisKeyUVDay.F(dayStart.Format("2006-01-02"))
			if n, err := cache.RDB.SCard(context.Background(), key).Result(); err == nil {
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
	if cache.RDB != nil {
		todayKey := common.RedisKeyUVDay.F(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Format("2006-01-02"))
		if n, err := cache.RDB.SCard(context.Background(), todayKey).Result(); err == nil {
			out.Totals.VisitorsToday = n
		}
	}
	return out, nil
}

