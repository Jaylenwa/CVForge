package db

import (
	"database/sql"
	"os"
	"path/filepath"

	"openresume/internal/infra/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
)

func RunMigrations(cfg config.Config, sqlDB *sql.DB) error {
	migrationsDir := filepath.Join("db", "migrations")
	if _, err := os.Stat(migrationsDir); err != nil {
		return nil
	}
	var (
		m   *migrate.Migrate
		err error
	)
	if cfg.SQLitePath != "" {
		driver, e := sqlite.WithInstance(sqlDB, &sqlite.Config{})
		if e != nil {
			return e
		}
		m, err = migrate.NewWithDatabaseInstance("file://"+migrationsDir, "sqlite3", driver)
	} else {
		driver, e := mysql.WithInstance(sqlDB, &mysql.Config{})
		if e != nil {
			return e
		}
		m, err = migrate.NewWithDatabaseInstance("file://"+migrationsDir, "mysql", driver)
	}
	if err != nil {
		return err
	}
	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
