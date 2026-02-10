package share

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"openresume/internal/infra/config"

	"github.com/golang-jwt/jwt/v5"
)

const shareTokenTyp = "share"

type shareTokenClaims struct {
	Typ  string `json:"typ"`
	Slug string `json:"slug"`
	V    string `json:"v"`
	jwt.RegisteredClaims
}

func shareSigningKey() []byte {
	return []byte("share:" + config.CF.JWTSecret)
}

func issueShareToken(sl ShareLink, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := shareTokenClaims{
		Typ:  shareTokenTyp,
		Slug: sl.Slug,
		V:    shareTokenVersion(sl),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tok.SignedString(shareSigningKey())
}

func shareTokenVersion(sl ShareLink) string {
	exp := ""
	if sl.ExpiresAt != nil {
		exp = fmt.Sprintf("%d", sl.ExpiresAt.Unix())
	}
	raw := strings.Join([]string{sl.Slug, sl.Password, exp}, "|")
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:8])
}

func validateShareToken(tokenStr string, sl ShareLink) error {
	if tokenStr == "" {
		return errors.New("missing token")
	}
	tok, err := jwt.ParseWithClaims(tokenStr, &shareTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing algorithm")
		}
		return shareSigningKey(), nil
	})
	if err != nil || tok == nil || !tok.Valid {
		return errors.New("invalid token")
	}
	claims, ok := tok.Claims.(*shareTokenClaims)
	if !ok || claims == nil {
		return errors.New("invalid claims")
	}
	if claims.Typ != shareTokenTyp {
		return errors.New("invalid typ")
	}
	if claims.Slug != sl.Slug {
		return errors.New("slug mismatch")
	}
	if claims.V == "" || claims.V != shareTokenVersion(sl) {
		return errors.New("token expired")
	}
	return nil
}
