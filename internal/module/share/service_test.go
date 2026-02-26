package share

import (
	"testing"
	"time"

	"cvforge/internal/infra/config"
	"cvforge/internal/infra/database"
	"cvforge/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(
		&models.User{},
		&models.Resume{},
		&models.ResumePersonal{},
		&models.ResumeTheme{},
		&models.ResumeSection{},
		&models.ResumeItem{},
		&models.ShareLink{},
	); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	database.DB = db
	config.CF = config.Config{JWTSecret: "test-secret"}
	return db
}

func TestSharePasswordAndExpiry(t *testing.T) {
	db := setupTestDB(t)
	email := "u@example.com"
	u := models.User{Email: &email, IsActive: true}
	if err := db.Create(&u).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	r := models.Resume{UserID: u.ID, Title: "V1", TemplateID: "TemplateClassic", Language: "en", LastModified: time.Now().UnixMilli()}
	if err := db.Create(&r).Error; err != nil {
		t.Fatalf("create resume: %v", err)
	}

	svc := NewService()
	sl, code, err := svc.PublishResumeForUser(u.ID, r.ID)
	if err != nil || code != 200 {
		t.Fatalf("publish: code=%d err=%v", code, err)
	}

	pass := "p@ssw0rd"
	_, code, err = svc.UpdateSettingsForUser(u.ID, r.ID, UpdateSettingsInput{
		IsPublic: ptrBool(true),
		Password: &pass,
	})
	if err != nil || code != 200 {
		t.Fatalf("update settings: code=%d err=%v", code, err)
	}

	_, code, err = svc.GetPublicPayload(sl.Slug, "")
	if code != 401 {
		t.Fatalf("expected password required, got code=%d err=%v", code, err)
	}

	_, code, err = svc.AuthenticatePublic(sl.Slug, "bad")
	if code != 401 {
		t.Fatalf("expected invalid password, got code=%d err=%v", code, err)
	}
	shareTok, code, err := svc.AuthenticatePublic(sl.Slug, pass)
	if err != nil || code != 200 || shareTok == "" {
		t.Fatalf("auth ok: code=%d err=%v token=%q", code, err, shareTok)
	}

	val, code, err := svc.GetPublicPayload(sl.Slug, shareTok)
	if err != nil || code != 200 {
		t.Fatalf("get public: code=%d err=%v", code, err)
	}
	if val == "" {
		t.Fatalf("expected payload")
	}

	if err := db.Model(&models.Resume{}).Where("id = ?", r.ID).Update("title", "V2").Error; err != nil {
		t.Fatalf("update resume: %v", err)
	}
	val2, code, err := svc.GetPublicPayload(sl.Slug, shareTok)
	if err != nil || code != 200 {
		t.Fatalf("get public2: code=%d err=%v", code, err)
	}
	if val2 == val {
		t.Fatalf("expected payload to change when resume updates")
	}

	exp := time.Now().Add(-time.Minute)
	_, code, err = svc.UpdateSettingsForUser(u.ID, r.ID, UpdateSettingsInput{ExpiresAt: &exp})
	if err != nil || code != 200 {
		t.Fatalf("set expires: code=%d err=%v", code, err)
	}
	_, code, _ = svc.GetPublicPayload(sl.Slug, shareTok)
	if code != 410 {
		t.Fatalf("expected expired 410, got %d", code)
	}
}

func ptrBool(v bool) *bool { return &v }
