package config

import (
	"context"
	"openresume/internal/models"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	return &Service{db: db, rdb: rdb}
}

func (s *Service) Get(key string) string {
	// Try cache first
	ctx := context.Background()
	val, err := s.rdb.Get(ctx, "sysconfig:"+key).Result()
	if err == nil {
		return val
	}

	// Try DB
	var cfg models.SystemConfig
	if err := s.db.Where("key = ?", key).First(&cfg).Error; err == nil {
		// Cache it
		s.rdb.Set(ctx, "sysconfig:"+key, cfg.Value, 24*time.Hour)
		return cfg.Value
	}

	return ""
}

func (s *Service) GetWithDefault(key, def string) string {
	val := s.Get(key)
	if val == "" {
		return def
	}
	return val
}

func (s *Service) GetBool(key string, def bool) bool {
	val := s.Get(key)
	if val == "" {
		return def
	}
	b, err := strconv.ParseBool(val)
	if err != nil {
		return def
	}
	return b
}

func (s *Service) Set(key, value, description, typeName string) error {
	var cfg models.SystemConfig
	err := s.db.Where("key = ?", key).First(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cfg = models.SystemConfig{
				Key:         key,
				Value:       value,
				Description: description,
				Type:        typeName,
			}
			if err := s.db.Create(&cfg).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		cfg.Value = value
		if description != "" {
			cfg.Description = description
		}
		if typeName != "" {
			cfg.Type = typeName
		}
		if err := s.db.Save(&cfg).Error; err != nil {
			return err
		}
	}

	// Invalidate cache
	s.rdb.Del(context.Background(), "sysconfig:"+key)
	return nil
}

func (s *Service) GetAll() ([]models.SystemConfig, error) {
	var configs []models.SystemConfig
	if err := s.db.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

type PublicConfig struct {
	EnableEmailVerification bool `json:"enableEmailVerification"`
	EnableWeChatLogin       bool `json:"enableWeChatLogin"`
	EnableGithubLogin       bool `json:"enableGithubLogin"`
	WeChatAppID             string `json:"weChatAppID"`
	GithubClientID          string `json:"githubClientID"`
}

func (s *Service) GetPublicConfig() PublicConfig {
	return PublicConfig{
		EnableEmailVerification: s.GetBool("enable_email_verification", false),
		EnableWeChatLogin:       s.GetBool("feature_wechat_login", true),
		EnableGithubLogin:       s.GetBool("feature_github_login", true),
		WeChatAppID:             s.Get("wechat_app_id"),
		GithubClientID:          s.Get("github_client_id"),
	}
}
