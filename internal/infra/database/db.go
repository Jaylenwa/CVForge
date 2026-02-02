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
	if err = dropTagColumns(); err != nil {
		return nil, err
	}

	logger.WithCtx(nil).Info("mysql connected")
	return DB, nil
}

func dropTagColumns() error {
	type target struct {
		table  string
		column string
	}
	targets := []target{
		{table: "template", column: "tags"},
		{table: "job_role", column: "tags"},
		{table: "content_preset", column: "tags"},
	}
	for _, t := range targets {
		if DB.Migrator().HasColumn(t.table, t.column) {
			if err := DB.Migrator().DropColumn(t.table, t.column); err != nil {
				return err
			}
		}
	}
	return nil
}

func autoMigrate() error {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.JobCategory{},
		&models.JobRole{},
		&models.ContentPreset{},
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
