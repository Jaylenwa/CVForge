package storage

import (
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/database"
)

func Init(cfg config.Config) error {
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
