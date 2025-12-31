package config

import (
	"os"
)

var CF Config

type Config struct {
	Port          string
	MySQLDSN      string
	SQLitePath    string
	RedisAddr     string
	RedisPassword string
	JWTSecret     string
}

func Load() Config {
	CF = Config{
		Port:          getenv("PORT", "8080"),
		MySQLDSN:      getenv("DB_DSN", ""),
		SQLitePath:    getenv("SQLITE_PATH", "openresume.db"),
		RedisAddr:     getenv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getenv("REDIS_PASSWORD", ""),
		JWTSecret:     getenv("JWT_SECRET", "abcdefghijklmnopqrstuvwxyz"),
	}
	return CF
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
