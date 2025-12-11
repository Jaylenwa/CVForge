package db

import (
	"database/sql"
	"log"

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
	if err = autoMigrate(); err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	log.Println("mysql connected")
	return sqlDB, nil
}

func Gorm() *gorm.DB { return g }

func autoMigrate() error {
	return g.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.Resume{},
		&models.ResumeSection{},
		&models.ResumeItem{},
		&models.ShareLink{},
	)
}
