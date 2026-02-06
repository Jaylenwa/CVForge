package storage

import (
	"context"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/database"
	"openresume/internal/models"
	svcConfig "openresume/internal/module/config"
	"openresume/internal/module/seed"
	tplmod "openresume/internal/module/template"
)

func Init() error {
	var err error
	defer func() {
		if err == nil {
			svcConfig.NewService().EnsureDefaults()
		}
	}()
	cfg := config.CF
	_, err = database.InitDB(cfg)
	if err != nil {
		return err
	}

	_, err = cache.InitRedis(cfg)
	if err != nil {
		return err
	}

	if e := autoSeedIfEmpty(); e != nil {
		return e
	}
	return nil
}

func autoSeedIfEmpty() error {
	var catCount int64
	var tplCount int64
	_ = database.DB.Model(&models.JobCategory{}).Count(&catCount).Error
	_ = database.DB.Model(&models.Template{}).Count(&tplCount).Error
	if catCount == 0 && tplCount == 0 {
		if s, e := seed.LoadDefaultSeed(); e == nil {
			if _, e2 := seed.Import(context.Background(), database.DB, s); e2 != nil {
				return e2
			}
		} else {
			return e
		}
		if items, e := seed.LoadDefaultTemplateItems(); e == nil && len(items) > 0 {
			if e2 := tplmod.NewService().Seed(items); e2 != nil {
				return e2
			}
		} else if e != nil {
			return e
		}
	}
	return nil
}
