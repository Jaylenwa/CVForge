package cache

import (
	"context"

	"openresume/config"

	"github.com/redis/go-redis/v9"
)

func InitRedis(cfg config.Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	_ = rdb.Ping(context.Background()).Err()
	return rdb
}
