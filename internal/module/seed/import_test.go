package seed

import (
	"context"
	"testing"

	"openresume/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestImport_DefaultSeed_Idempotent(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"))
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(
		&models.JobCategory{},
		&models.JobCategoryI18n{},
		&models.JobRole{},
		&models.JobRoleI18n{},
		&models.ContentPreset{},
		&models.ContentPresetI18n{},
	); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	if err := db.Exec("CREATE UNIQUE INDEX IF NOT EXISTS uniq_content_preset_language ON content_preset_i18n (content_preset_id, language)").Error; err != nil {
		t.Fatalf("ensure uniq index: %v", err)
	}

	s, err := LoadDefaultSeed()
	if err != nil {
		t.Fatalf("load default seed: %v", err)
	}

	if _, err := Import(context.Background(), db, s); err != nil {
		t.Fatalf("import: %v", err)
	}

	var c1, r1, p1, p1i int64
	if err := db.Model(&models.JobCategory{}).Count(&c1).Error; err != nil {
		t.Fatalf("count categories: %v", err)
	}
	if err := db.Model(&models.JobRole{}).Count(&r1).Error; err != nil {
		t.Fatalf("count roles: %v", err)
	}
	if err := db.Model(&models.ContentPreset{}).Count(&p1).Error; err != nil {
		t.Fatalf("count presets: %v", err)
	}
	if err := db.Model(&models.ContentPresetI18n{}).Count(&p1i).Error; err != nil {
		t.Fatalf("count preset i18n: %v", err)
	}

	if _, err := Import(context.Background(), db, s); err != nil {
		t.Fatalf("import again: %v", err)
	}

	var c2, r2, p2, p2i int64
	if err := db.Model(&models.JobCategory{}).Count(&c2).Error; err != nil {
		t.Fatalf("count categories again: %v", err)
	}
	if err := db.Model(&models.JobRole{}).Count(&r2).Error; err != nil {
		t.Fatalf("count roles again: %v", err)
	}
	if err := db.Model(&models.ContentPreset{}).Count(&p2).Error; err != nil {
		t.Fatalf("count presets again: %v", err)
	}
	if err := db.Model(&models.ContentPresetI18n{}).Count(&p2i).Error; err != nil {
		t.Fatalf("count preset i18n again: %v", err)
	}

	if c1 != c2 || r1 != r2 || p1 != p2 || p1i != p2i {
		t.Fatalf("not idempotent: categories %d->%d roles %d->%d presets %d->%d preset_i18n %d->%d", c1, c2, r1, r2, p1, p2, p1i, p2i)
	}

	if c1 != int64(len(s.Categories)) {
		t.Fatalf("unexpected category count: got %d want %d", c1, len(s.Categories))
	}
	if r1 != int64(len(s.Roles)) {
		t.Fatalf("unexpected role count: got %d want %d", r1, len(s.Roles))
	}
	if p1 != int64(len(s.Presets)) {
		t.Fatalf("unexpected preset base count: got %d want %d", p1, len(s.Presets))
	}
	if p1i != int64(len(s.Presets)) {
		t.Fatalf("unexpected preset i18n count: got %d want %d", p1i, len(s.Presets))
	}
}

