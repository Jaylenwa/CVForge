package share

import (
	"context"
	"encoding/json"
	"time"

	"openresume/internal/common"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repo
	db   *gorm.DB
	rdb  *redis.Client
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{repo: NewRepo(db), db: db, rdb: rdb}
}

func (s *Service) PublishResumeForUser(uid uint, externalID string) (ShareLink, int, error) {
	var res Resume
	if err := s.db.Where("external_id = ?", externalID).First(&res).Error; err != nil {
		return ShareLink{}, 404, err
	}
	if res.UserID != uid {
		return ShareLink{}, 403, gorm.ErrInvalidData
	}
	var sl ShareLink
	if err := s.db.Where("resume_id = ?", res.ID).First(&sl).Error; err != nil {
		sl = ShareLink{ResumeID: res.ID, Slug: uuid.NewString()[:8], IsPublic: true}
		_ = s.db.Create(&sl).Error
	} else {
		sl.IsPublic = true
		_ = s.db.Save(&sl).Error
	}
	return sl, 200, nil
}

func (s *Service) GetPublicPayload(slug string) (string, int, error) {
	cacheKey := common.RedisKeyPublicResume.F(slug)
	if s.rdb != nil {
		if val, err := s.rdb.Get(context.Background(), cacheKey).Result(); err == nil {
			_ = s.rdb.Incr(context.Background(), common.RedisKeyViews.F(slug))
			return val, 200, nil
		}
	}
	sl, err := s.repo.FindPublicBySlug(slug)
	if err != nil {
		return "", 404, err
	}
	var res Resume
	if err := s.db.Where("id = ?", sl.ResumeID).Preload("Sections.Items").First(&res).Error; err != nil {
		return "", 404, err
	}
	payloadObj := gin.H{
		"Title":        res.Title,
		"TemplateID":   res.TemplateID,
		"ThemeColor":   res.ThemeColor,
		"ThemeFont":    res.ThemeFont,
		"ThemeSpacing": res.ThemeSpacing,
		"FullName":     res.FullName,
		"Email":        res.Email,
		"Phone":        res.Phone,
		"AvatarURL":    res.AvatarURL,
		"Sections":     res.Sections,
	}
	b, _ := json.Marshal(payloadObj)
	val := string(b)
	if s.rdb != nil {
		_ = s.rdb.Set(context.Background(), cacheKey, val, 10*time.Minute).Err()
		_ = s.rdb.Incr(context.Background(), common.RedisKeyViews.F(slug))
	}
	return val, 200, nil
}
