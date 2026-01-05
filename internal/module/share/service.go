package share

import (
	"context"
	"encoding/json"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/database"
	resmod "openresume/internal/module/resume"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) PublishResumeForUser(uid uint, externalID string) (ShareLink, int, error) {
	var res Resume
	if err := database.DB.Where("external_id = ?", externalID).First(&res).Error; err != nil {
		return ShareLink{}, 404, err
	}
	if res.UserID != uid {
		return ShareLink{}, 403, gorm.ErrInvalidData
	}
	var sl ShareLink
	if err := database.DB.Where("resume_id = ?", res.ID).First(&sl).Error; err != nil {
		sl = ShareLink{ResumeID: res.ID, UserID: uid, Slug: uuid.NewString()[:8], IsPublic: true}
		_ = database.DB.Create(&sl).Error
	} else {
		sl.IsPublic = true
		if sl.UserID == 0 {
			sl.UserID = uid
		}
		_ = database.DB.Save(&sl).Error
	}
	return sl, 200, nil
}

func (s *Service) GetPublicPayload(slug string) (string, int, error) {
	cacheKey := common.RedisKeyPublicResume.F(slug)
	if cache.RDB != nil {
		if val, err := cache.RDB.Get(context.Background(), cacheKey).Result(); err == nil {
			_ = cache.RDB.Incr(context.Background(), common.RedisKeyViews.F(slug))
			return val, 200, nil
		}
	}
	sl, err := s.repo.FindPublicBySlug(slug)
	if err != nil {
		return "", 404, err
	}
	var res Resume
	if err := database.DB.Where("id = ?", sl.ResumeID).Preload("Personal").Preload("Job").Preload("Theme").Preload("Sections.Items").First(&res).Error; err != nil {
		return "", 404, err
	}
	payload := resmod.ToDTO(res)
	b, _ := json.Marshal(payload)
	val := string(b)
	if cache.RDB != nil {
		_ = cache.RDB.Set(context.Background(), cacheKey, val, 10*time.Minute).Err()
		_ = cache.RDB.Incr(context.Background(), common.RedisKeyViews.F(slug))
	}
	return val, 200, nil
}
