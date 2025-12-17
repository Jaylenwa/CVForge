package storage

import (
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/db"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Init(cfg config.Config) (*gorm.DB, *redis.Client, error) {
	db, err := db.InitMySQL(cfg)
	if err != nil {
		return nil, nil, err
	}
	rdb, err := cache.InitRedis(cfg)
	if err != nil {
		return nil, nil, err
	}
	return db, rdb, nil
}
