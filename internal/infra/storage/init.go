package storage

import (
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/database"
)

func Init() error {
	cfg := config.CF
	_, err := database.InitMySQL(cfg)
	if err != nil {
		return err
	}
	_, err = cache.InitRedis(cfg)
	if err != nil {
		return err
	}
	return nil
}
