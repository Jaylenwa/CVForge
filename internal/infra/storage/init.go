package storage

import (
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/database"
	svcConfig "openresume/internal/module/config"
)

func Init() error {
	var err error
	defer func() {
		if err == nil {
			svcConfig.NewService().EnsureDefaults()
		}
	}()
	cfg := config.CF
	_, err = database.InitMySQL(cfg)
	if err != nil {
		return err
	}
	_, err = cache.InitRedis(cfg)
	if err != nil {
		return err
	}
	return nil
}
