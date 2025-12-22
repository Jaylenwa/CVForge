package config

import (
	"context"
	"openresume/internal/infra/config"
	"openresume/internal/models"
	"os"
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

// EnsureDefaults inserts a set of default configuration keys into DB if they don't exist.
// Values are initialized from environment-backed Config where applicable.
func (s *Service) EnsureDefaults(cfg config.Config) error {
	type def struct {
		Key         string
		Value       string
		Description string
		Type        string
	}
	defaults := []def{
		// Registration
		{"enable_email_verification", "false", "Enable email verification during registration", "bool"},
		// SMTP
		{"smtp_host", os.Getenv("SMTP_HOST"), "SMTP host", "string"},
		{"smtp_port", os.Getenv("SMTP_PORT"), "SMTP port", "string"},
		{"smtp_username", os.Getenv("SMTP_USERNAME"), "SMTP username", "string"},
		{"smtp_password", os.Getenv("SMTP_PASSWORD"), "SMTP password", "string"},
		{"smtp_from_name", os.Getenv("SMTP_FROM_NAME"), "SMTP from name", "string"},
		// WeChat OAuth
		{"wechat_login", "false", "Enable WeChat login", "bool"},
		{"wechat_app_id", os.Getenv("WECHAT_APP_ID"), "WeChat AppID", "string"},
		{"wechat_app_secret", os.Getenv("WECHAT_APP_SECRET"), "WeChat App Secret", "string"},
		{"wechat_redirect_uri", os.Getenv("WECHAT_REDIRECT_URI"), "WeChat Redirect URI", "string"},
		// GitHub OAuth
		{"github_login", "false", "Enable GitHub login", "bool"},
		{"github_client_id", os.Getenv("GITHUB_CLIENT_ID"), "GitHub Client ID", "string"},
		{"github_client_secret", os.Getenv("GITHUB_CLIENT_SECRET"), "GitHub Client Secret", "string"},
		{"github_redirect_uri", os.Getenv("GITHUB_REDIRECT_URI"), "GitHub Redirect URI", "string"},
	}
	for _, d := range defaults {
		var existing models.Config
		if err := s.db.Where("key = ?", d.Key).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				_ = s.db.Create(&models.Config{
					Key:         d.Key,
					Value:       d.Value,
					Description: d.Description,
					Type:        d.Type,
				}).Error
				// best-effort; continue on error to attempt other inserts
				continue
			}
			// if other DB error, abort early
			return err
		}
	}
	return nil
}

func (s *Service) Get(key string) string {
	// Try cache first
	ctx := context.Background()
	val, err := s.rdb.Get(ctx, "sysconfig:"+key).Result()
	if err == nil {
		return val
	}

	// Try DB
	var cfg models.Config
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
	var cfg models.Config
	err := s.db.Where("key = ?", key).First(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cfg = models.Config{
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

func (s *Service) GetAll() ([]models.Config, error) {
	var configs []models.Config
	if err := s.db.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

type PublicConfig struct {
	EnableEmailVerification bool `json:"enableEmailVerification"`
	EnableWeChatLogin       bool `json:"enableWeChatLogin"`
	EnableGithubLogin       bool `json:"enableGithubLogin"`
}

func (s *Service) GetPublicConfig() PublicConfig {
	// Read from DB once and map required keys
	var (
		enableEmailVerification bool
		enableWeChatLogin       bool
		enableGithubLogin       bool
	)
	configs, err := s.GetAll()
	if err == nil {
		for _, cfg := range configs {
			switch cfg.Key {
			case "enable_email_verification":
				if b, err := strconv.ParseBool(cfg.Value); err == nil {
					enableEmailVerification = b
				}
			case "wechat_login":
				if b, err := strconv.ParseBool(cfg.Value); err == nil {
					enableWeChatLogin = b
				}
			case "github_login":
				if b, err := strconv.ParseBool(cfg.Value); err == nil {
					enableGithubLogin = b
				}
			}
		}
	}
	return PublicConfig{
		EnableEmailVerification: enableEmailVerification,
		EnableWeChatLogin:       enableWeChatLogin,
		EnableGithubLogin:       enableGithubLogin,
	}
}
