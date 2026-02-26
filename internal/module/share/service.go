package share

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	"cvforge/internal/infra/database"
	resmod "cvforge/internal/module/resume"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) PublishResumeForUser(uid uint, id uint) (ShareLink, int, error) {
	var res Resume
	if err := database.DB.Where("id = ?", id).First(&res).Error; err != nil {
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

var (
	ErrPasswordRequired = errors.New("password required")
	ErrExpired          = errors.New("expired")
)

func (s *Service) GetPublicPayload(slug string, shareToken string) (string, int, error) {
	sl, err := s.repo.FindBySlug(slug)
	if err != nil {
		return "", 404, err
	}
	if !sl.IsPublic {
		return "", 404, gorm.ErrRecordNotFound
	}
	if sl.ExpiresAt != nil && time.Now().After(*sl.ExpiresAt) {
		return "", 410, ErrExpired
	}
	if sl.Password != "" {
		if err := validateShareToken(shareToken, sl); err != nil {
			return "", 401, ErrPasswordRequired
		}
	}

	cacheKey := common.RedisKeyPublicResume.F(slug)
	if cache.RDB != nil {
		if val, err := cache.RDB.Get(context.Background(), cacheKey).Result(); err == nil {
			s.trackAccess(context.Background(), sl)
			return val, 200, nil
		}
	}
	var res Resume
	if err := database.DB.Where("id = ?", sl.ResumeID).Preload("Personal").Preload("Theme").Preload("Sections.Items").First(&res).Error; err != nil {
		return "", 404, err
	}
	payload := resmod.ToDTO(res)
	b, _ := json.Marshal(payload)
	val := string(b)
	if cache.RDB != nil {
		_ = cache.RDB.Set(context.Background(), cacheKey, val, 10*time.Minute).Err()
	}
	s.trackAccess(context.Background(), sl)
	return val, 200, nil
}

func (s *Service) trackAccess(ctx context.Context, sl ShareLink) {
	now := time.Now()
	if cache.RDB != nil {
		_ = cache.RDB.Incr(ctx, common.RedisKeyViews.F(sl.Slug))
		_ = cache.RDB.Set(ctx, common.RedisKeyShareLastAccess.F(sl.Slug), now.UnixMilli(), 30*24*time.Hour).Err()
		ok := cache.RDB.SetNX(ctx, "share:stats:flush:"+sl.Slug, "1", 30*time.Second).Val()
		if !ok {
			return
		}
		views := cache.RDB.Get(ctx, common.RedisKeyViews.F(sl.Slug)).Val()
		var cnt uint64
		if views != "" {
			var n uint64
			_, _ = fmt.Sscanf(views, "%d", &n)
			cnt = n
		}
		_ = database.DB.Model(&ShareLink{}).Where("id = ?", sl.ID).Updates(map[string]any{
			"views":          cnt,
			"last_access_at": now,
		}).Error
		return
	}
	_ = database.DB.Model(&ShareLink{}).Where("id = ?", sl.ID).Updates(map[string]any{
		"views":          gorm.Expr("views + 1"),
		"last_access_at": now,
	}).Error
}

func (s *Service) AuthenticatePublic(slug string, password string) (string, int, error) {
	sl, err := s.repo.FindBySlug(slug)
	if err != nil {
		return "", 404, err
	}
	if !sl.IsPublic {
		return "", 404, gorm.ErrRecordNotFound
	}
	if sl.ExpiresAt != nil && time.Now().After(*sl.ExpiresAt) {
		return "", 410, ErrExpired
	}
	if sl.Password == "" {
		return "", 400, errors.New("password not set")
	}
	if subtle.ConstantTimeCompare([]byte(sl.Password), []byte(password)) != 1 {
		return "", 401, errors.New("invalid password")
	}
	tok, err := issueShareToken(sl, 30*time.Minute)
	if err != nil {
		return "", 500, err
	}
	return tok, 200, nil
}

type UpdateSettingsInput struct {
	IsPublic       *bool
	Password       *string
	ExpiresAt      *time.Time
	ClearExpiresAt bool
}

func (s *Service) GetSettingsForUser(uid uint, resumeID uint) (ShareLink, int, error) {
	var res Resume
	if err := database.DB.Where("id = ?", resumeID).First(&res).Error; err != nil {
		return ShareLink{}, 404, err
	}
	if res.UserID != uid {
		return ShareLink{}, 403, gorm.ErrInvalidData
	}
	sl, err := s.repo.FindByResumeID(res.ID)
	if err != nil {
		return ShareLink{}, 404, err
	}
	if sl.UserID == 0 {
		sl.UserID = uid
		_ = s.repo.Save(&sl)
	}
	return sl, 200, nil
}

func (s *Service) UpdateSettingsForUser(uid uint, resumeID uint, in UpdateSettingsInput) (ShareLink, int, error) {
	var res Resume
	if err := database.DB.Where("id = ?", resumeID).First(&res).Error; err != nil {
		return ShareLink{}, 404, err
	}
	if res.UserID != uid {
		return ShareLink{}, 403, gorm.ErrInvalidData
	}
	sl, err := s.repo.FindByResumeID(res.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sl = ShareLink{ResumeID: res.ID, UserID: uid, Slug: uuid.NewString()[:8], IsPublic: false}
			if err := s.repo.Create(&sl); err != nil {
				return ShareLink{}, 500, err
			}
		} else {
			return ShareLink{}, 500, err
		}
	}
	changed := false
	if in.IsPublic != nil {
		sl.IsPublic = *in.IsPublic
		changed = true
	}
	if in.ClearExpiresAt {
		sl.ExpiresAt = nil
		changed = true
	} else if in.ExpiresAt != nil {
		sl.ExpiresAt = in.ExpiresAt
		changed = true
	}
	if in.Password != nil {
		p := *in.Password
		if p == "" {
			sl.Password = ""
		} else {
			sl.Password = p
		}
		changed = true
	}
	if sl.UserID == 0 {
		sl.UserID = uid
		changed = true
	}
	if changed {
		if err := s.repo.Save(&sl); err != nil {
			return ShareLink{}, 500, err
		}
		if cache.RDB != nil {
			_ = cache.RDB.Del(context.Background(), common.RedisKeyPublicResume.F(sl.Slug)).Err()
		}
	}
	return sl, 200, nil
}
