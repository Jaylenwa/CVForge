package config

import (
	"os"
)

type Config struct {
	Port                string
	MySQLDSN            string
	SQLitePath          string
	RedisAddr           string
	RedisPassword       string
	JWTSecret           string
	CORSOrigins         string
	GeminiAPIKey        string
	UploadBackend       string
	S3Bucket            string
	S3Region            string
	S3Endpoint          string
	S3AccessKey         string
	S3SecretKey         string
	FrontendBaseURL     string
	ChromeJSONURL       string
	SMTPHost            string
	SMTPPort            string
	SMTPUsername        string
	SMTPPassword        string
	SMTPFromName        string
	WeChatAppID         string
	WeChatAppSecret     string
	WeChatRedirectURI   string
	GithubClientID      string
	GithubClientSecret  string
	GithubRedirectURI   string
	OAuthAllowedOrigins string
	FeatureWeChatLogin  string
	FeatureGithubLogin  string
}

func Load() Config {
	return Config{
		Port:                getenv("PORT", "8080"),
		MySQLDSN:            getenv("DB_DSN", "root:password@tcp(127.0.0.1:3306)/openresume?charset=utf8mb4&parseTime=True&loc=Local"),
		SQLitePath:          getenv("SQLITE_PATH", ""),
		RedisAddr:           getenv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:       getenv("REDIS_PASSWORD", ""),
		JWTSecret:           getenv("JWT_SECRET", "devsecret"),
		CORSOrigins:         getenv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173,http://182.254.166.74:8889/api/v1/auth/github/callback"),
		GeminiAPIKey:        getenv("GEMINI_API_KEY", ""),
		UploadBackend:       getenv("UPLOAD_BACKEND", "local"),
		S3Bucket:            getenv("S3_BUCKET", ""),
		S3Region:            getenv("S3_REGION", ""),
		S3Endpoint:          getenv("S3_ENDPOINT", ""),
		S3AccessKey:         getenv("S3_ACCESS_KEY", ""),
		S3SecretKey:         getenv("S3_SECRET_KEY", ""),
		FrontendBaseURL:     getenv("FRONTEND_BASE_URL", "http://localhost:3000"),
		ChromeJSONURL:       getenv("CHROME_JSON_URL", "http://localhost:9222/json/version"),
		SMTPHost:            getenv("SMTP_HOST", ""),
		SMTPPort:            getenv("SMTP_PORT", ""),
		SMTPUsername:        getenv("SMTP_USERNAME", ""),
		SMTPPassword:        getenv("SMTP_PASSWORD", ""),
		SMTPFromName:        getenv("SMTP_FROM_NAME", "OpenResume"),
		WeChatAppID:         getenv("WECHAT_APP_ID", ""),
		WeChatAppSecret:     getenv("WECHAT_APP_SECRET", ""),
		WeChatRedirectURI:   getenv("WECHAT_REDIRECT_URI", ""),
		GithubClientID:      getenv("GITHUB_CLIENT_ID", ""),
		GithubClientSecret:  getenv("GITHUB_CLIENT_SECRET", ""),
		GithubRedirectURI:   getenv("GITHUB_REDIRECT_URI", "http://182.254.166.74:8889/api/v1/auth/github/callback"),
		OAuthAllowedOrigins: getenv("OAUTH_ALLOWED_ORIGINS", getenv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173,http://127.0.0.1:3000,http://127.0.0.1:5173,http://182.254.166.74:8889")),
		FeatureWeChatLogin:  getenv("FEATURE_WECHAT_LOGIN", "on"),
		FeatureGithubLogin:  getenv("FEATURE_GITHUB_LOGIN", "on"),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
