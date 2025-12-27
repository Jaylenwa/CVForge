package pdf

import (
	"context"
	"encoding/json"
	"fmt"
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
	type verInfo struct {
		WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
	}
	var v verInfo
	resp, err := http.Get(cfg.ChromeJSONURL)
	if err != nil {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil || v.WebSocketDebuggerUrl == "" {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, fmt.Errorf("ws empty")
	}
	u, _ := url.Parse(v.WebSocketDebuggerUrl)
	// 强制替换 host
	u.Host = "chrome:3000"
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
	var pdf []byte
	err = chromedp.Run(ctx,
		chromedp.Navigate(cfg.FrontendBaseURL),
		chromedp.ActionFunc(func(ctx context.Context) error {
			js := fmt.Sprintf("localStorage.setItem('token', %s)", strconv.Quote(token))
			return chromedp.Evaluate(js, nil).Do(ctx)
		}),
		chromedp.Navigate(dest),
		emulation.SetDeviceMetricsOverride(1200, 1800, 1, false),
		chromedp.WaitVisible(`#resume-export-root`, chromedp.ByID),
		chromedp.Evaluate(`(function(){var s=document.createElement('style');s.innerHTML='@media print{body *{visibility:hidden} #resume-export-root,#resume-export-root *{visibility:visible} #resume-export-root{position:absolute;left:0;top:0}}';document.head.appendChild(s)})()`, nil),
		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, perr := page.PrintToPDF().WithPrintBackground(true).
				WithDisplayHeaderFooter(false).WithPaperWidth(8.27).
				WithPaperHeight(11.69).WithMarginTop(0.5).
				WithMarginBottom(0.5).WithMarginLeft(0.5).WithMarginRight(0.5).Do(ctx)
			if perr != nil {
				return perr
			}
			pdf = buf
			return nil
		}),
	)
	if err != nil {
		s.cbFail(common.CBCircuitPDF)
		return nil, 503, err
	}
	s.cbReset(common.CBCircuitPDF)
	return pdf, 200, nil
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
	resp, err := http.Get(cfg.ChromeJSONURL)
	if err != nil {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil || v.WebSocketDebuggerUrl == "" {
		s.cbFail(common.CBCircuitImage)
		return nil, 503, fmt.Errorf("ws empty")
	}
	allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(context.Background(), v.WebSocketDebuggerUrl)
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
