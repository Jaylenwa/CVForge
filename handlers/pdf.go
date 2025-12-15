package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"openresume/config"
	"openresume/models"

	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterPDFRoutes(r *gin.RouterGroup, db *gorm.DB, auth gin.HandlerFunc) {
	r.POST("/resumes/:id/pdf", auth, func(c *gin.Context) {
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).Preload("Sections.Items").First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		cfg := config.Load()
		type verInfo struct {
			WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
		}
		var v verInfo
		resp, err := http.Get(cfg.ChromeJSONURL)
		if err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "pdf service unavailable"})
			return
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil || v.WebSocketDebuggerUrl == "" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "pdf service unavailable"})
			return
		}
		allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(context.Background(), v.WebSocketDebuggerUrl)
		defer cancelAlloc()
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()
		if cfg.FrontendBaseURL == "" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "pdf service unavailable"})
			return
		}
		var pdf []byte
		dest := cfg.FrontendBaseURL + "/#/print?id=" + c.Param("id")
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
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
			c.JSON(http.StatusNotImplemented, gin.H{"error": "pdf service unavailable"})
			return
		}
		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", "attachment; filename=resume.pdf")
		c.Writer.Write(pdf)
	})

	r.POST("/resumes/:id/image", auth, func(c *gin.Context) {
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).Preload("Sections.Items").First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		cfg := config.Load()
		type verInfo struct {
			WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
		}
		var v verInfo
		resp, err := http.Get(cfg.ChromeJSONURL)
		if err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "image service unavailable"})
			return
		}
		defer resp.Body.Close()
		if err := json.NewDecoder(resp.Body).Decode(&v); err != nil || v.WebSocketDebuggerUrl == "" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "image service unavailable"})
			return
		}
		allocCtx, cancelAlloc := chromedp.NewRemoteAllocator(context.Background(), v.WebSocketDebuggerUrl)
		defer cancelAlloc()
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()
		if cfg.FrontendBaseURL == "" {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "image service unavailable"})
			return
		}
		var png []byte
		dest := cfg.FrontendBaseURL + "/#/print?id=" + c.Param("id")
		authHeader := c.GetHeader("Authorization")
		token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
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
				var perr error
				root, perr := dom.GetDocument().Do(ctx)
				if perr != nil {
					return perr
				}
				nid, perr := dom.QuerySelector(root.NodeID, "#resume-export-root").Do(ctx)
				if perr != nil {
					return perr
				}
				bm, perr := dom.GetBoxModel().WithNodeID(nid).Do(ctx)
				if perr != nil {
					return perr
				}
				c := bm.Content
				if len(c) < 8 {
					return nil
				}
				vx := float64(c[0])
				vy := float64(c[1])
				vw := float64(c[2] - c[0])
				vh := float64(c[7] - c[1])
				buf, perr := page.CaptureScreenshot().WithClip(&page.Viewport{X: vx, Y: vy, Width: vw, Height: vh, Scale: 1}).Do(ctx)
				if perr != nil {
					return perr
				}
				png = buf
				return nil
			}),
		)
		if err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "image service unavailable"})
			return
		}
		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", "attachment; filename=resume.png")
		c.Writer.Write(png)
	})
}
