package database

import (
	"os"

	"openresume/internal/infra/config"
	"openresume/internal/models"
	"openresume/internal/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg config.Config) (*gorm.DB, error) {
	var err error
	if cfg.MySQLDSN != "" {
		DB, err = gorm.Open(mysql.Open(cfg.MySQLDSN))
	} else {
		DB, err = gorm.Open(sqlite.Open(cfg.SQLitePath))
	}
	if os.Getenv("GORM_LOG") == "DEV" {
		DB.Logger = gormLogger.Default.LogMode(gormLogger.Info)
	}
	if err != nil {
		return nil, err
	}

	if err = autoMigrate(); err != nil {
		return nil, err
	}
	if err = ensureContentPresetI18nIndex(); err != nil {
		return nil, err
	}
	// removed legacy backfill

	logger.WithCtx(nil).Info("mysql connected")
	return DB, nil
}

func ensureContentPresetI18nIndex() error {
	if DB == nil {
		return nil
	}
	dialect := DB.Dialector.Name()
	switch dialect {
	case "sqlite":
		_ = DB.Exec("DROP INDEX IF EXISTS uniq_content_preset_id").Error
		_ = DB.Exec("CREATE UNIQUE INDEX IF NOT EXISTS uniq_content_preset_language ON content_preset_i18n (content_preset_id, language)").Error
	default:
		_ = DB.Exec("DROP INDEX uniq_content_preset_id ON content_preset_i18n").Error
		_ = DB.Exec("CREATE UNIQUE INDEX uniq_content_preset_language ON content_preset_i18n (content_preset_id, language)").Error
	}
	return nil
}

func autoMigrate() error {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.TemplateI18n{},
		&models.JobCategory{},
		&models.JobCategoryI18n{},
		&models.JobRole{},
		&models.JobRoleI18n{},
		&models.ContentPreset{},
		&models.ContentPresetI18n{},
		&models.RoleTemplateUsage{},
		&models.Resume{},
		&models.ResumePersonal{},
		&models.ResumeTheme{},
		&models.ResumeSection{},
		&models.ResumeItem{},
		&models.ShareLink{},
		&models.AuditLog{},
		&models.OAuthAccount{},
		&models.Config{},
	); err != nil {
		return err
	}
	return nil
}

// legacy backfill removed
