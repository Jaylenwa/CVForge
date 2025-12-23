package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestDailyUVIdentity_UID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{Header: make(http.Header)}
	c.Set("uid", uint(123))
	id := dailyUVIdentity(c)
	if id != "u:123" {
		t.Fatalf("expected u:123, got %s", id)
	}
}

func TestDailyUVIdentity_IP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	req := &http.Request{Header: make(http.Header)}
	req.RemoteAddr = "9.8.7.6:12345"
	c.Request = req
	id := dailyUVIdentity(c)
	if id != "i:9.8.7.6" {
		t.Fatalf("expected i:9.8.7.6, got %s", id)
	}
}
