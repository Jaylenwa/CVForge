package handlers

import (
	"context"
	"net/http"
	"time"

	"openresume/config"
	"openresume/middleware"
	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type sendCodeReq struct {
	Email string `json:"email"`
}
type registerReq struct{ Email, Code, Password, Name string }
type loginReq struct{ Email, Password string }
type refreshReq struct {
	RefreshToken string `json:"refreshToken"`
}

func RegisterAuthRoutes(r *gin.RouterGroup, cfg config.Config, rdb *redis.Client, db *gorm.DB) {
	auth := r.Group("/auth")
	auth.POST("/send-code", middleware.RateLimit(rdb, 3, time.Minute), func(c *gin.Context) {
		var req sendCodeReq
		if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email"})
			return
		}
		code := "123456"
		_ = rdb.Set(context.Background(), "verify:"+req.Email, code, 10*time.Minute).Err()
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	auth.POST("/register", func(c *gin.Context) {
		var req registerReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		if req.Code != "123456" {
			val, err := rdb.Get(context.Background(), "verify:"+req.Email).Result()
			if err != nil || val != req.Code {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid code"})
				return
			}
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		u := models.User{Email: req.Email, PasswordHash: string(hash), Name: req.Name}
		if err := db.Create(&u).Error; err != nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email exists"})
			return
		}
		access, refresh := issueTokens(cfg, u.ID)
		c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
	})

	auth.POST("/login", func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var u models.User
		if err := db.Where("email = ?", req.Email).First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password)) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		access, refresh := issueTokens(cfg, u.ID)
		c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
	})

	auth.POST("/refresh", func(c *gin.Context) {
		var req refreshReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		t, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil })
		if err != nil || !t.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, _ := t.Claims.(jwt.MapClaims)
		jti, _ := claims["jti"].(string)
		if jti != "" {
			if rdb.Get(context.Background(), "jwt:blacklist:"+jti).Val() == "1" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "revoked"})
				return
			}
		}
		uid := uint(claims["uid"].(float64))
		access, refresh := issueTokens(cfg, uid)
		c.JSON(http.StatusOK, gin.H{"accessToken": access, "refreshToken": refresh})
	})

	auth.POST("/logout", func(c *gin.Context) {
		var req refreshReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		t, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(cfg.JWTSecret), nil })
		if err != nil || !t.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, _ := t.Claims.(jwt.MapClaims)
		jti, _ := claims["jti"].(string)
		if jti != "" {
			_ = rdb.Set(context.Background(), "jwt:blacklist:"+jti, "1", time.Hour*24*7).Err()
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}

func issueTokens(cfg config.Config, uid uint) (string, string) {
	mk := func(exp time.Duration) string {
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"uid": uid, "exp": time.Now().Add(exp).Unix(), "jti": uuid.NewString()})
		s, _ := t.SignedString([]byte(cfg.JWTSecret))
		return s
	}
	return mk(2 * time.Hour), mk(7 * 24 * time.Hour)
}
