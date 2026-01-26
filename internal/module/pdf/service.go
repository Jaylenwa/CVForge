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
	"openresume/internal/infra/database"
	conf "openresume/internal/module/config"
	"openresume/internal/pkg/storage"

	"github.com/gin-gonic/gin"
)

type Service struct {
	sysConfig *conf.Service
}

func NewService() *Service {
	return &Service{sysConfig: conf.NewService()}
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

func (s *Service) tryMarkConsumed(ctx context.Context, key string) bool {
	if cache.RDB == nil {
		return false
	}
	return cache.RDB.SetNX(ctx, key, "1", 10*time.Minute).Val()
}

func (s *Service) scheduleDelete(ctx context.Context, jobID string) {
	time.Sleep(5 * time.Second)
	s.deleteJobFile(ctx, jobID)
}

func (s *Service) deleteJobFile(ctx context.Context, jobID string) {
	resultURL := cache.RDB.Get(ctx, jobKey(jobID)+":result").Val()
	if resultURL == "" {
		return
	}
	cfg := s.sysConfig.GetStorageSettings()
	if up, err := storage.NewFromSettings(cfg); err == nil {
		_ = up.Delete(ctx, resultURL)
	}
	_ = cache.RDB.Del(ctx, jobKey(jobID)+":result").Err()
	_ = cache.RDB.Del(ctx, jobKey(jobID)+":status").Err()
	_ = cache.RDB.Del(ctx, jobKey(jobID)+":data").Err()
	_ = cache.RDB.Del(ctx, jobKey(jobID)+":filename").Err()
	_ = cache.RDB.Del(ctx, jobKey(jobID)+":consumed").Err()
}

func (s *Service) GeneratePDFWithToken(id uint, token string) ([]byte, int, error) {
	if s.cbOpen(common.CBCircuitPDF) {
		return nil, 503, fmt.Errorf("cb")
	}
	var res Resume
	if err := database.DB.Where("id = ?", id).Preload("Sections.Items").First(&res).Error; err != nil {
		return nil, 404, err
	}
	if s.sysConfig.Get(string(common.ConfigKeyFrontendBaseURL)) == "" {
		return nil, 503, fmt.Errorf("fe empty")
	}
	dest := s.sysConfig.Get(string(common.ConfigKeyFrontendBaseURL)) + "/#/print?id=" + fmt.Sprintf("%d", id)
	if token == "" {
		token = ""
	}
	reqBody := PDFRequest{
		URL: dest,
		Options: map[string]any{
			"format":            "A4",
			"printBackground":   true,
			"preferCSSPageSize": true,
			"margin": map[string]string{
				"top":    "0px",
				"bottom": "0px",
				"left":   "0px",
				"right":  "0px",
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

func (s *Service) GenerateImage(c *gin.Context, id uint) ([]byte, int, error) {
	if s.cbOpen(common.CBCircuitImage) {
		return nil, 503, fmt.Errorf("cb")
	}
	var res Resume
	if err := database.DB.Where("id = ?", id).Preload("Sections.Items").First(&res).Error; err != nil {
		return nil, 404, err
	}
	if s.sysConfig.Get(string(common.ConfigKeyFrontendBaseURL)) == "" {
		return nil, 503, fmt.Errorf("fe empty")
	}
	dest := s.sysConfig.Get(string(common.ConfigKeyFrontendBaseURL)) + "/#/print?id=" + fmt.Sprintf("%d", id)
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
