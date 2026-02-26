package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	"cvforge/internal/infra/config"
	"cvforge/internal/infra/database"
	"cvforge/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
				return nil, errors.New("invalid signing algorithm")
			}
			return []byte(config.CF.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
			return
		}
		if expVal, ok := claims["exp"]; ok {
			switch t := expVal.(type) {
			case float64:
				if time.Now().Unix() > int64(t) {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
					return
				}
			case int64:
				if time.Now().Unix() > t {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token expired"})
					return
				}
			}
		}
		if cache.RDB != nil {
			if jti, _ := claims["jti"].(string); jti != "" {
				if cache.RDB.Get(context.Background(), common.RedisKeyJWTBlacklist.F(jti)).Val() == "1" {
					c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
					return
				}
			}
		}
		uid := uint(claims["uid"].(float64))
		var tokenVer int = 1
		if v, ok := claims["ver"]; ok {
			switch t := v.(type) {
			case float64:
				if int(t) > 0 {
					tokenVer = int(t)
				}
			case int64:
				if int(t) > 0 {
					tokenVer = int(t)
				}
			case int:
				if t > 0 {
					tokenVer = t
				}
			}
		}
		active := ""
		verCached := ""
		if cache.RDB != nil {
			active, _ = cache.RDB.Get(context.Background(), common.RedisKeyUserActive.F(fmt.Sprintf("%d", uid))).Result()
			verCached, _ = cache.RDB.Get(context.Background(), common.RedisKeyUserTokenVersion.F(fmt.Sprintf("%d", uid))).Result()
		}
		if active == "0" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
		if verCached != "" && verCached != "1" {
			if verCached == "inc" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}
		}
		if active == "" {
			var u models.User
			if err := database.DB.First(&u, uid).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				return
			}
			if !u.IsActive {
				if cache.RDB != nil {
					_ = cache.RDB.Set(context.Background(), common.RedisKeyUserActive.F(fmt.Sprintf("%d", uid)), "0", time.Minute).Err()
				}
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
			if u.TokenVersion > 0 && tokenVer != u.TokenVersion {
				if cache.RDB != nil {
					_ = cache.RDB.Set(context.Background(), common.RedisKeyUserTokenVersion.F(fmt.Sprintf("%d", uid)), "inc", time.Minute).Err()
				}
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
				return
			}
			if cache.RDB != nil {
				_ = cache.RDB.Set(context.Background(), common.RedisKeyUserActive.F(fmt.Sprintf("%d", uid)), "1", time.Minute).Err()
				_ = cache.RDB.Set(context.Background(), common.RedisKeyUserTokenVersion.F(fmt.Sprintf("%d", uid)), "1", time.Minute).Err()
			}
		}
		c.Set("uid", uid)
		c.Next()
	}
}
