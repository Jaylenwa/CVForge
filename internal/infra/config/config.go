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
	FrontendBaseURL string
	ChromeJSONURL   string
}

func Load() Config {
	return Config{
		Port:            getenv("PORT", "8080"),
		MySQLDSN:        getenv("DB_DSN", "root:password@tcp(127.0.0.1:3306)/openresume?charset=utf8mb4&parseTime=True&loc=Local"),
		SQLitePath:      getenv("SQLITE_PATH", ""),
		RedisAddr:       getenv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword:   getenv("REDIS_PASSWORD", ""),
		JWTSecret:       getenv("JWT_SECRET", "devsecret"),
		CORSOrigins:     getenv("CORS_ORIGINS", "http://localhost:3000,http://127.0.0.1:3000"),
		FrontendBaseURL: getenv("FRONTEND_BASE_URL", "http://localhost:3000"),
		ChromeJSONURL:   getenv("CHROME_JSON_URL", ""),
	}
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
