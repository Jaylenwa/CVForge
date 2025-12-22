package config

import (
	"context"
	"openresume/internal/infra/config"
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
		{"smtp_host", cfg.SMTPHost, "SMTP host", "string"},
		{"smtp_port", cfg.SMTPPort, "SMTP port", "string"},
		{"smtp_username", cfg.SMTPUsername, "SMTP username", "string"},
		{"smtp_password", cfg.SMTPPassword, "SMTP password", "string"},
		{"smtp_from_name", cfg.SMTPFromName, "SMTP from name", "string"},
		// WeChat OAuth
		{"feature_wechat_login", "false", "Enable WeChat login", "bool"},
		{"wechat_app_id", cfg.WeChatAppID, "WeChat AppID", "string"},
		{"wechat_app_secret", cfg.WeChatAppSecret, "WeChat App Secret", "string"},
		{"wechat_redirect_uri", cfg.WeChatRedirectURI, "WeChat Redirect URI", "string"},
		// GitHub OAuth
		{"feature_github_login", "false", "Enable GitHub login", "bool"},
		{"github_client_id", cfg.GithubClientID, "GitHub Client ID", "string"},
		{"github_client_secret", cfg.GithubClientSecret, "GitHub Client Secret", "string"},
		{"github_redirect_uri", cfg.GithubRedirectURI, "GitHub Redirect URI", "string"},
	}
	for _, d := range defaults {
		var existing models.SystemConfig
		if err := s.db.Where("key = ?", d.Key).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				_ = s.db.Create(&models.SystemConfig{
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
	EnableEmailVerification bool   `json:"enableEmailVerification"`
	EnableWeChatLogin       bool   `json:"enableWeChatLogin"`
	EnableGithubLogin       bool   `json:"enableGithubLogin"`
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
