package cache

import (
	"context"

	"cvforge/internal/infra/config"

	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

func InitRedis(cfg config.Config) (*redis.Client, error) {
	RDB = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	err := RDB.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return RDB, nil
}
