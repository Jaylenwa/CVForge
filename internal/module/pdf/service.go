package pdf

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/database"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
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
	cfg := config.Load()
	if cfg.FrontendBaseURL == "" {
		return nil, 503, fmt.Errorf("fe empty")
	}
	dest := "http://frontend/#/print?id=" + externalID
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	type PDFRequest struct {
		URL              string         `json:"url"`
		HTML             string         `json:"html,omitempty"`
		Options          map[string]any `json:"options,omitempty"`
		EmulateMediaType string         `json:"emulateMediaType,omitempty"`
		GotoOptions      struct {
			Referer        string   `json:"referer,omitempty"`
			ReferrerPolicy string   `json:"referrerPolicy,omitempty"`
			Timeout        int      `json:"timeout,omitempty"`
			WaitUntil      []string `json:"waitUntil,omitempty"`
		} `json:"gotoOptions,omitempty"`
		WaitForSelector struct {
			Hidden   bool   `json:"hidden,omitempty"`
			Selector string `json:"selector,omitempty"`
			Timeout  int    `json:"timeout,omitempty"`
			Visible  bool   `json:"visible,omitempty"`
		} `json:"waitForSelector,omitempty"`
		SetExtraHTTPHeaders map[string]string `json:"setExtraHTTPHeaders,omitempty"`
		BestAttempt         bool              `json:"bestAttempt,omitempty"`
	}
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
	pdfAPI := cfg.ChromeAPIURL + "/pdf"
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
	cfg := config.Load()
	type verInfo struct {
		WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
	}
	var v verInfo
	resp, err := http.Get(cfg.ChromeAPIURL)
	if err != nil {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil || v.WebSocketDebuggerUrl == "" {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, fmt.Errorf("ws empty")
	}
	u, _ := url.Parse(v.WebSocketDebuggerUrl)
	u.Host = "vfoy.cn:3000"
	wsURL := u.String()
	allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(context.Background(), wsURL)
	defer cancelAlloc()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	if cfg.FrontendBaseURL == "" {
		return nil, 503, fmt.Errorf("fe empty")
	}
	dest := cfg.FrontendBaseURL + "/#/print?id=" + externalID
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	var png []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(cfg.FrontendBaseURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			js := fmt.Sprintf("localStorage.setItem('token', %s)", strconv.Quote(token))
			return chromedp.Evaluate(js, nil).Do(ctx)
		}),
		chromedp.Navigate(dest),
		emulation.SetDeviceMetricsOverride(1200, 1800, 1, false),
		chromedp.WaitVisible(`#resume-export-root`, chromedp.ByID),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, perr := page.CaptureScreenshot().Do(ctx)
			if perr != nil {
				return perr
			}
			png = buf
			return nil
		}),
	)
	if err != nil {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, err
	}
	s.cbReset(common.CBCircuitImage)
	return png, 200, nil
}
