package database

import (
	"testing"

	"openresume/internal/infra/config"
)

func TestAutoMigrateSQLite(t *testing.T) {
	cfg := config.Config{
		SQLitePath: t.TempDir() + "/test.db",
	}
	db, err := InitMySQL(cfg)
	if err != nil {
		t.Fatalf("InitMySQL failed: %v", err)
	}
	sqlDB, err := db.DB()
	if err == nil {
		_ = sqlDB.Close()
	}
}

