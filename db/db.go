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
	if err := g.AutoMigrate(
		&models.User{},
		&models.Template{},
		&models.Resume{},
		&models.ResumeSection{},
		&models.ResumeItem{},
		&models.ShareLink{},
	); err != nil {
		return err
	}
	// Drop deprecated columns
	if g.Migrator().HasColumn(&models.ResumeItem{}, "location") {
		if err := g.Migrator().DropColumn(&models.ResumeItem{}, "location"); err != nil {
			return err
		}
	}
	if g.Migrator().HasColumn(&models.Resume{}, "address") {
		if err := g.Migrator().DropColumn(&models.Resume{}, "address"); err != nil {
			return err
		}
	}
	return nil
}
