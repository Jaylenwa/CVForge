package pdf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/database"
	conf "openresume/internal/module/config"

	"github.com/gin-gonic/gin"
)

type Service struct {
	cfg       config.Config
	sysConfig *conf.Service
}

func NewService(cfg config.Config, sys *conf.Service) *Service {
	return &Service{cfg: cfg, sysConfig: sys}
}

func (s *Service) cbOpen(svc string) bool {
	if cache.RDB == nil {
		return false
	}
	return cache.RDB.Get(context.Background(), common.RedisKeyCircuitBreaker.F(svc)).Val() == common.CBCircuitOpen
}
func (s *Service) cbFail(svc string) {
	if cache.RDB == nil {
		return
	}
	cnt, _ := cache.RDB.Incr(context.Background(), common.RedisKeyCircuitBreakerFail.F(svc)).Result()
	if cnt == 1 {
		_ = cache.RDB.Expire(context.Background(), common.RedisKeyCircuitBreakerFail.F(svc), time.Minute).Err()
	}
	if cnt >= 3 {
		_ = cache.RDB.Set(context.Background(), common.RedisKeyCircuitBreaker.F(svc), common.CBCircuitOpen, time.Minute).Err()
	}
}
func (s *Service) cbReset(svc string) {
	if cache.RDB == nil {
		return
	}
	_ = cache.RDB.Del(context.Background(), common.RedisKeyCircuitBreakerFail.F(svc)).Err()
	_ = cache.RDB.Del(context.Background(), common.RedisKeyCircuitBreaker.F(svc)).Err()
}

func (s *Service) GeneratePDF(c *gin.Context, externalID string) ([]byte, int, error) {
	if s.cbOpen(common.CBCircuitPDF) {
		return nil, 503, fmt.Errorf("cb")
	}
	var res Resume
	if err := database.DB.Where("external_id = ?", externalID).Preload("Sections.Items").First(&res).Error; err != nil {
		return nil, 404, err
	}
	if s.cfg.FrontendBaseURL == "" {
		return nil, 503, fmt.Errorf("fe empty")
	}
	dest := "http://frontend/#/print?id=" + externalID
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	reqBody := PDFRequest{
		URL: dest,
		Options: map[string]any{
			"format":          "A4",
			"printBackground": true,
			"margin": map[string]any{
				"top":    "20mm",
				"bottom": "20mm",
				"left":   "15mm",
				"right":  "15mm",
			},
		},
		EmulateMediaType: "print",
		GotoOptions: struct {
			Referer        string   `json:"referer,omitempty"`
			ReferrerPolicy string   `json:"referrerPolicy,omitempty"`
			Timeout        int      `json:"timeout,omitempty"`
			WaitUntil      []string `json:"waitUntil,omitempty"`
		}{
			Timeout:   60000,
			WaitUntil: []string{"networkidle0"},
		},
		WaitForSelector: struct {
			Hidden   bool   `json:"hidden,omitempty"`
			Selector string `json:"selector,omitempty"`
			Timeout  int    `json:"timeout,omitempty"`
			Visible  bool   `json:"visible,omitempty"`
		}{
			Selector: "#resume-export-root",
			Visible:  true,
			Timeout:  60000,
		},
		SetExtraHTTPHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
		BestAttempt: true,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, err
	}
	chrome := s.sysConfig.Get(string(common.ConfigKeyChromeAPIURL))
	if chrome == "" {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, fmt.Errorf("chrome empty")
	}
	pdfAPI := chrome + "/pdf"
	resp, err := http.Post(pdfAPI, "application/json", bytes.NewReader(payload))
	if err != nil {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, fmt.Errorf("call pdf api failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, fmt.Errorf("pdf api status: %d", resp.StatusCode)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, err
	}
	s.cbReset(common.CBCircuitPDF)
	return buf, 200, nil
}

func (s *Service) GenerateImage(c *gin.Context, externalID string) ([]byte, int, error) {
	if s.cbOpen(common.CBCircuitImage) {
		return nil, 503, fmt.Errorf("cb")
	}
	var res Resume
	if err := database.DB.Where("external_id = ?", externalID).Preload("Sections.Items").First(&res).Error; err != nil {
		return nil, 404, err
	}
	if s.cfg.FrontendBaseURL == "" {
		return nil, 503, fmt.Errorf("fe empty")
	}
	dest := "http://frontend/#/print?id=" + externalID
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	reqBody := ScreenshotRequest{
		URL: dest,
		Options: map[string]any{
			"type":           "png",
			"omitBackground": false,
			"fullPage":       false,
		},
		EmulateMediaType: "print",
		GotoOptions: struct {
			Referer        string   `json:"referer,omitempty"`
			ReferrerPolicy string   `json:"referrerPolicy,omitempty"`
			Timeout        int      `json:"timeout,omitempty"`
			WaitUntil      []string `json:"waitUntil,omitempty"`
		}{
			Timeout:   60000,
			WaitUntil: []string{"networkidle0"},
		},
		WaitForSelector: struct {
			Hidden   bool   `json:"hidden,omitempty"`
			Selector string `json:"selector,omitempty"`
			Timeout  int    `json:"timeout,omitempty"`
			Visible  bool   `json:"visible,omitempty"`
		}{
			Selector: "#resume-export-root",
			Visible:  true,
			Timeout:  60000,
		},
		Selector:             "#resume-export-root",
		SetJavaScriptEnabled: true,
		Viewport: struct {
			Width             int     `json:"width,omitempty"`
			Height            int     `json:"height,omitempty"`
			DeviceScaleFactor float64 `json:"deviceScaleFactor,omitempty"`
			IsMobile          bool    `json:"isMobile,omitempty"`
			IsLandscape       bool    `json:"isLandscape,omitempty"`
			HasTouch          bool    `json:"hasTouch,omitempty"`
		}{
			Width:             1200,
			Height:            1800,
			DeviceScaleFactor: 1,
			IsMobile:          false,
			IsLandscape:       false,
			HasTouch:          false,
		},
		SetExtraHTTPHeaders: map[string]string{
			"Authorization": "Bearer " + token,
		},
		BestAttempt: true,
	}
	payload, err := json.Marshal(reqBody)
	if err != nil {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, err
	}
	chrome := s.sysConfig.Get(string(common.ConfigKeyChromeAPIURL))
	if chrome == "" {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, fmt.Errorf("chrome empty")
	}
	imgAPI := chrome + "/screenshot"
	resp, err := http.Post(imgAPI, "application/json", bytes.NewReader(payload))
	if err != nil {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, fmt.Errorf("call screenshot api failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, fmt.Errorf("screenshot api status: %d", resp.StatusCode)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, err
	}
	s.cbReset(common.CBCircuitImage)
	return buf, 200, nil
}
