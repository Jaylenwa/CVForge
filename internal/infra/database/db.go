package database

import (
	"log"
	"os"

	"openresume/internal/infra/config"
	"openresume/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(cfg config.Config) (*gorm.DB, error) {
	var err error
	if cfg.MySQLDSN != "" {
		DB, err = gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	} else {
		DB, err = gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	}
	if err != nil {
		return nil, err
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, err
	}
	// prefer versioned migrations if present; otherwise fall back to AutoMigrate for initial schema
	if _, e := os.Stat("db/migrations"); e == nil {
		if err = RunMigrations(cfg, sqlDB); err != nil {
			return nil, err
		}
	} else {
		if err = autoMigrate(); err != nil {
			return nil, err
		}
	}
	log.Println("mysql connected")
	return DB, nil
}

func autoMigrate() error {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Template{},
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
