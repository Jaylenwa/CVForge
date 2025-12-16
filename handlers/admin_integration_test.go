package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"openresume/config"
	"openresume/db"
	"openresume/middleware"
	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	cfg := config.Config{
		SQLitePath: ":memory:",
		MySQLDSN:   "",
	}
	_, err := db.InitMySQL(cfg)
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	return db.Gorm()
}

func makeToken(secret string, uid uint, exp time.Duration) string {
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid": uid,
		"exp": time.Now().Add(exp).Unix(),
	})
	s, _ := tok.SignedString([]byte(secret))
	return s
}

func TestAdminUsersList(t *testing.T) {
	os.Setenv("GIN_MODE", "test")
	gin.SetMode(gin.TestMode)
	dbx := setupTestDB(t)
	admin := models.User{Email: "admin@example.com", PasswordHash: "", Name: "Admin", Role: "admin", IsActive: true}
	user := models.User{Email: "user@example.com", PasswordHash: "", Name: "User", Role: "user", IsActive: true}
	if err := dbx.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}
	if err := dbx.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	cfg := config.Config{JWTSecret: "devsecret"}
	r := gin.New()
	api := r.Group("/api/v1")
	RegisterAdminUserRoutes(api, dbx, middleware.Auth(cfg), middleware.RequireRole("admin"))
	token := makeToken(cfg.JWTSecret, admin.ID, time.Hour)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 got %d body=%s", w.Code, w.Body.String())
	}
}

func TestAdminUsersListForbidden(t *testing.T) {
	os.Setenv("GIN_MODE", "test")
	gin.SetMode(gin.TestMode)
	dbx := setupTestDB(t)
	admin := models.User{Email: "admin2@example.com", PasswordHash: "", Name: "Admin2", Role: "admin", IsActive: true}
	user := models.User{Email: "user2@example.com", PasswordHash: "", Name: "User2", Role: "user", IsActive: true}
	if err := dbx.Create(&admin).Error; err != nil {
		t.Fatalf("create admin: %v", err)
	}
	if err := dbx.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	cfg := config.Config{JWTSecret: "devsecret"}
	r := gin.New()
	api := r.Group("/api/v1")
	RegisterAdminUserRoutes(api, dbx, middleware.Auth(cfg), middleware.RequireRole("admin"))
	token := makeToken(cfg.JWTSecret, user.ID, time.Hour)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users?page=1&pageSize=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d body=%s", w.Code, w.Body.String())
	}
}
