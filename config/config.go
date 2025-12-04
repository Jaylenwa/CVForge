package config

import (
	"os"
)

type Config struct {
	Port            string
	MySQLDSN        string
	SQLitePath      string
	RedisAddr       string
	RedisPassword   string
	JWTSecret       string
	CORSOrigins     string
	GeminiAPIKey    string
	UploadBackend   string
	S3Bucket        string
	S3Region        string
	S3Endpoint      string
	S3AccessKey     string
	S3SecretKey     string
	FrontendBaseURL string
}

func Load() Config {
	return Config{
		Port:            getenv("PORT", "8080"),
		MySQLDSN:        getenv("DB_DSN", "root:password@tcp(127.0.0.1:3306)/openresume?charset=utf8mb4&parseTime=True&loc=Local"),
		SQLitePath:      getenv("SQLITE_PATH", ""),
		RedisAddr:       getenv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:   getenv("REDIS_PASSWORD", ""),
		JWTSecret:       getenv("JWT_SECRET", "devsecret"),
		CORSOrigins:     getenv("CORS_ORIGINS", "http://localhost:3000,http://localhost:5173"),
		GeminiAPIKey:    getenv("GEMINI_API_KEY", ""),
		UploadBackend:   getenv("UPLOAD_BACKEND", "local"),
		S3Bucket:        getenv("S3_BUCKET", ""),
		S3Region:        getenv("S3_REGION", ""),
		S3Endpoint:      getenv("S3_ENDPOINT", ""),
		S3AccessKey:     getenv("S3_ACCESS_KEY", ""),
		S3SecretKey:     getenv("S3_SECRET_KEY", ""),
		FrontendBaseURL: getenv("FRONTEND_BASE_URL", ""),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
