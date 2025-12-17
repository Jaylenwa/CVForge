package cache

import (
	"context"

	"openresume/internal/infra/config"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

func InitRedis(cfg config.Config) (*redis.Client, error) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       0,
	})
	err := rdb.Ping(context.Background()).Err()
	if err != nil {
		return nil, err
	}
	return rdb, nil
}

func Redis() *redis.Client { return rdb }
