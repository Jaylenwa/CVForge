package config

import (
	"context"
	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/database"
	"openresume/internal/models"
	"openresume/internal/pkg/storage"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

// EnsureDefaults inserts a set of default configuration keys into DB if they don't exist.
// Values are initialized from environment-backed Config where applicable.
func (s *Service) EnsureDefaults() error {
	type def struct {
		Key         string
		Value       string
		Description string
		Type        string
	}
	defaults := []def{
		// Registration
		{string(common.ConfigKeyEnableEmailVerification), "false", "Enable email verification during registration", "bool"},
		{string(common.ConfigKeyEnabledPricingPage), "false", "Enable Pricing page", "bool"},
		{string(common.ConfigKeySMTPHost), os.Getenv("SMTP_HOST"), "SMTP host", "string"},
		{string(common.ConfigKeySMTPPort), os.Getenv("SMTP_PORT"), "SMTP port", "string"},
		{string(common.ConfigKeySMTPUsername), os.Getenv("SMTP_USERNAME"), "SMTP username", "string"},
		{string(common.ConfigKeySMTPPassword), os.Getenv("SMTP_PASSWORD"), "SMTP password", "string"},
		{string(common.ConfigKeySMTPFromName), os.Getenv("SMTP_FROM_NAME"), "SMTP from name", "string"},
		// WeChat OAuth
		{string(common.ConfigKeyEnabledWechatLogin), "false", "Enable WeChat login", "bool"},
		{string(common.ConfigKeyWeChatAppID), os.Getenv("WECHAT_APP_ID"), "WeChat AppID", "string"},
		{string(common.ConfigKeyWeChatAppSecret), os.Getenv("WECHAT_APP_SECRET"), "WeChat App Secret", "string"},
		{string(common.ConfigKeyWeChatRedirectURI), os.Getenv("WECHAT_REDIRECT_URI"), "WeChat Redirect URI", "string"},
		// WeChat Official Account (MP) QR login
		{string(common.ConfigKeyEnabledWechatMPLogin), "false", "Enable WeChat MP QR login", "bool"},
		{string(common.ConfigKeyWeChatMPAppID), os.Getenv("WECHAT_MP_APP_ID"), "WeChat MP AppID", "string"},
		{string(common.ConfigKeyWeChatMPAppSecret), os.Getenv("WECHAT_MP_APP_SECRET"), "WeChat MP App Secret", "string"},
		{string(common.ConfigKeyWeChatMPToken), os.Getenv("WECHAT_MP_TOKEN"), "WeChat MP callback token", "string"},
		{string(common.ConfigKeyWeChatMPAESKey), os.Getenv("WECHAT_MP_AES_KEY"), "WeChat MP callback AES key (optional)", "string"},
		// GitHub OAuth
		{string(common.ConfigKeyEnabledGithubLogin), "false", "Enable GitHub login", "bool"},
		{string(common.ConfigKeyGithubClientID), os.Getenv("GITHUB_CLIENT_ID"), "GitHub Client ID", "string"},
		{string(common.ConfigKeyGithubClientSecret), os.Getenv("GITHUB_CLIENT_SECRET"), "GitHub Client Secret", "string"},
		{string(common.ConfigKeyGithubRedirectURI), os.Getenv("GITHUB_REDIRECT_URI"), "GitHub Redirect URI", "string"},
		// Storage
		{string(common.ConfigKeyEnabledStorageS3), "false", "Enable S3 storage", "bool"},
		{string(common.ConfigKeyStorageS3Bucket), os.Getenv("S3_BUCKET"), "S3 bucket", "string"},
		{string(common.ConfigKeyStorageS3Region), os.Getenv("S3_REGION"), "S3 region", "string"},
		{string(common.ConfigKeyStorageS3Endpoint), os.Getenv("S3_ENDPOINT"), "S3 endpoint (optional)", "string"},
		{string(common.ConfigKeyStorageS3AccessKey), os.Getenv("S3_ACCESS_KEY"), "S3 access key (optional)", "string"},
		{string(common.ConfigKeyStorageS3SecretKey), os.Getenv("S3_SECRET_KEY"), "S3 secret key (optional)", "string"},
		{string(common.ConfigKeyChromeAPIURL), os.Getenv("CHROME_API_URL"), "Chrome API base URL", "string"},
		// Security
		{string(common.ConfigKeyCORSOrigins), os.Getenv("CORS_ORIGINS"), "CORS allowed origins (comma separated)", "string"},
		{string(common.ConfigKeyFrontendBaseURL), os.Getenv("FRONTEND_BASE_URL"), "Frontend base URL", "string"},
		// PDF worker controls
		{"pdf_worker_concurrency", "2", "PDF worker concurrency", "int"},
		{"pdf_queue_max", "1000", "PDF queue max length", "int"},
	}
	for _, d := range defaults {
		var existing models.Config
		if err := database.DB.Where("config_key = ?", d.Key).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				_ = database.DB.Create(&models.Config{
					ConfigKey:   d.Key,
					ConfigValue: d.Value,
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
	ctx := context.Background()
	if cache.RDB != nil {
		if val, err := cache.RDB.Get(ctx, common.RedisKeySysConfig.F(key)).Result(); err == nil {
			return val
		}
	}

	var cfg models.Config
	if err := database.DB.Where("config_key = ?", key).First(&cfg).Error; err == nil {
		if cache.RDB != nil {
			cache.RDB.Set(ctx, common.RedisKeySysConfig.F(key), cfg.ConfigValue, 24*time.Hour)
		}
		return cfg.ConfigValue
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
	err := database.DB.Where("config_key = ?", key).First(&cfg).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			cfg = models.Config{
				ConfigKey:   key,
				ConfigValue: value,
				Description: description,
				Type:        typeName,
			}
			if err := database.DB.Create(&cfg).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		cfg.ConfigValue = value
		if description != "" {
			cfg.Description = description
		}
		if typeName != "" {
			cfg.Type = typeName
		}
		if err := database.DB.Save(&cfg).Error; err != nil {
			return err
		}
	}

	// Invalidate cache
	if cache.RDB != nil {
		cache.RDB.Del(context.Background(), common.RedisKeySysConfig.F(key))
	}
	return nil
}

func (s *Service) GetAll() ([]models.Config, error) {
	var configs []models.Config
	if err := database.DB.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

type PublicConfig struct {
	EnableEmailVerification bool `json:"enableEmailVerification"`
	EnableWeChatLogin       bool `json:"enableWeChatLogin"`
	EnableWeChatMPLogin     bool `json:"enableWeChatMPLogin"`
	EnableGithubLogin       bool `json:"enableGithubLogin"`
	EnablePricingPage       bool `json:"enablePricingPage"`
}

func (s *Service) GetPublicConfig() PublicConfig {
	// Read from DB once and map required keys
	var (
		enableEmailVerification bool
		enableWeChatLogin       bool
		enableWeChatMPLogin     bool
		enableGithubLogin       bool
		enablePricingPage       bool
	)
	configs, err := s.GetAll()
	if err == nil {
		for _, cfg := range configs {
			switch cfg.ConfigKey {
			case string(common.ConfigKeyEnableEmailVerification):
				if b, err := strconv.ParseBool(cfg.ConfigValue); err == nil {
					enableEmailVerification = b
				}
			case string(common.ConfigKeyEnabledWechatLogin):
				if b, err := strconv.ParseBool(cfg.ConfigValue); err == nil {
					enableWeChatLogin = b
				}
			case string(common.ConfigKeyEnabledWechatMPLogin):
				if b, err := strconv.ParseBool(cfg.ConfigValue); err == nil {
					enableWeChatMPLogin = b
				}
			case string(common.ConfigKeyEnabledGithubLogin):
				if b, err := strconv.ParseBool(cfg.ConfigValue); err == nil {
					enableGithubLogin = b
				}
			case string(common.ConfigKeyEnabledPricingPage):
				if b, err := strconv.ParseBool(cfg.ConfigValue); err == nil {
					enablePricingPage = b
				}
			}
		}
	}
	return PublicConfig{
		EnableEmailVerification: enableEmailVerification,
		EnableWeChatLogin:       enableWeChatLogin,
		EnableWeChatMPLogin:     enableWeChatMPLogin,
		EnableGithubLogin:       enableGithubLogin,
		EnablePricingPage:       enablePricingPage,
	}
}

func (s *Service) GetStorageSettings() storage.Settings {
	return storage.Settings{
		Enabled:   s.GetBool(string(common.ConfigKeyEnabledStorageS3), false),
		Bucket:    s.Get(string(common.ConfigKeyStorageS3Bucket)),
		Region:    s.Get(string(common.ConfigKeyStorageS3Region)),
		Endpoint:  s.Get(string(common.ConfigKeyStorageS3Endpoint)),
		AccessKey: s.Get(string(common.ConfigKeyStorageS3AccessKey)),
		SecretKey: s.Get(string(common.ConfigKeyStorageS3SecretKey)),
	}
}
