package db

import (
	"database/sql"
	"log"
	"os"

	"openresume/config"
	"openresume/models"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var g *gorm.DB

func InitMySQL(cfg config.Config) (*sql.DB, error) {
	var (
		db  *gorm.DB
		err error
	)
	if cfg.SQLitePath != "" {
		db, err = gorm.Open(sqlite.Open(cfg.SQLitePath), &gorm.Config{})
	} else {
		db, err = gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	}
	if err != nil {
		return nil, err
	}
	g = db
	sqlDB, err := db.DB()
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
	return sqlDB, nil
}

func Gorm() *gorm.DB { return g }

func autoMigrate() error {
	if err := g.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.Resume{},
		&models.ResumeSection{},
		&models.ResumeItem{},
		&models.ShareLink{},
		&models.AuditLog{},
	); err != nil {
		return err
	}
	return nil
}
