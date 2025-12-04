package handlers

import (
	"context"
	"net/http"
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
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		var pdf []byte
		header := buildHeader(res)
		footer := buildFooter()
		var err error
		if cfg.FrontendBaseURL != "" {
			dest := cfg.FrontendBaseURL + "/#/print?id=" + c.Param("id")
			err = chromedp.Run(ctx,
				chromedp.Navigate(dest),
				emulation.SetDeviceMetricsOverride(1200, 1800, 1, false),
				// 仅打印目标元素
				chromedp.Evaluate(`(function(){var s=document.createElement('style');s.innerHTML='@media print{body *{visibility:hidden} #resume-export-root,#resume-export-root *{visibility:visible} #resume-export-root{position:absolute;left:0;top:0}}';document.head.appendChild(s)})()`, nil),
				chromedp.ActionFunc(func(ctx context.Context) error {
					buf, _, err := page.PrintToPDF().WithPrintBackground(true).WithDisplayHeaderFooter(false).WithPaperWidth(8.27).WithPaperHeight(11.69).WithMarginTop(0.5).WithMarginBottom(0.5).WithMarginLeft(0.5).WithMarginRight(0.5).Do(ctx)
					if err != nil {
						return err
					}
					pdf = buf
					return nil
				}),
			)
		} else {
			html := renderResumeHTML(res)
			url := "data:text/html," + urlEncode(html)
			err = chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.ActionFunc(func(ctx context.Context) error {
					buf, _, err := page.PrintToPDF().WithPrintBackground(true).WithDisplayHeaderFooter(true).WithHeaderTemplate(header).WithFooterTemplate(footer).WithPaperWidth(8.27).WithPaperHeight(11.69).WithMarginTop(0.5).WithMarginBottom(0.5).WithMarginLeft(0.5).WithMarginRight(0.5).Do(ctx)
					if err != nil {
						return err
					}
					pdf = buf
					return nil
				}),
			)
		}
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
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()
		var png []byte
		var err error
		if cfg.FrontendBaseURL != "" {
			dest := cfg.FrontendBaseURL + "/#/print?id=" + c.Param("id")
			err = chromedp.Run(ctx,
				chromedp.Navigate(dest),
				emulation.SetDeviceMetricsOverride(1200, 1800, 1, false),
				chromedp.ActionFunc(func(ctx context.Context) error {
					// 定位目标元素并按其边界截图
					root, err := dom.GetDocument().Do(ctx)
					if err != nil {
						return err
					}
					nid, err := dom.QuerySelector(root.NodeID, "#resume-export-root").Do(ctx)
					if err != nil {
						return err
					}
					bm, err := dom.GetBoxModel().WithNodeID(nid).Do(ctx)
					if err != nil {
						return err
					}
					c := bm.Content
					if len(c) < 8 {
						return nil
					}
					vx := float64(c[0])
					vy := float64(c[1])
					vw := float64(c[2] - c[0])
					vh := float64(c[7] - c[1])
					buf, err := page.CaptureScreenshot().WithClip(&page.Viewport{X: vx, Y: vy, Width: vw, Height: vh, Scale: 1}).Do(ctx)
					if err != nil {
						return err
					}
					png = buf
					return nil
				}),
			)
		} else {
			html := renderResumeHTML(res)
			url := "data:text/html," + urlEncode(html)
			err = chromedp.Run(ctx,
				chromedp.Navigate(url),
				chromedp.FullScreenshot(&png, 100),
			)
		}
		if err != nil {
			c.JSON(http.StatusNotImplemented, gin.H{"error": "image service unavailable"})
			return
		}
		c.Header("Content-Type", "image/png")
		c.Header("Content-Disposition", "attachment; filename=resume.png")
		c.Writer.Write(png)
	})
}

func renderResumeHTML(r models.Resume) string {
	accent := r.ThemeColor
	if accent == "" {
		accent = "#2563eb"
	}
	var b strings.Builder
	b.WriteString("<html><head><meta charset=\"utf-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><style>")
	b.WriteString("@page{size:A4;margin:20mm} body{font-family:'Inter',Arial,'Microsoft YaHei',sans-serif;color:#1f2937;background:#fff;} .container{padding:0 12px}")
	b.WriteString(".header{display:flex;align-items:flex-end;justify-content:space-between;border-bottom:4px solid " + accent + ";padding-bottom:8px;margin-bottom:16px}")
	b.WriteString(".name{font-size:26px;font-weight:800;letter-spacing:.5px;color:" + accent + "} .title{font-size:14px;color:#6b7280}")
	b.WriteString(".contact{margin-top:6px;color:#6b7280;font-size:12px}")
	b.WriteString(".section{margin-top:16px} .sec-title{font-size:12px;font-weight:700;text-transform:uppercase;letter-spacing:.14em;color:#374151;margin-bottom:8px}")
	b.WriteString(".item{margin-bottom:10px} .item-top{display:flex;justify-content:space-between;align-items:baseline} .item-title{font-weight:700;color:#111827}")
	b.WriteString(".item-meta{font-size:12px;color:" + accent + "} .desc{margin-top:6px;font-size:12px;line-height:1.7;white-space:pre-wrap;color:#374151}")
	b.WriteString(".badge{display:inline-block;margin:2px 6px 6px 0;padding:4px 8px;border:1px solid #e5e7eb;border-radius:10px;font-size:12px;color:#374151;background:#fff}")
	b.WriteString("</style></head><body><div class=container>")
	b.WriteString("<div class=header>")
	b.WriteString("<div><div class=name>" + escape(r.FullName) + "</div><div class=title>" + escape(r.TemplateID) + "</div>")
	b.WriteString("<div class=contact>")
	if r.Email != "" {
		b.WriteString(escape(r.Email))
	}
	if r.Phone != "" {
		b.WriteString(" • " + escape(r.Phone))
	}
	if r.Address != "" {
		b.WriteString(" • " + escape(r.Address))
	}
	if r.Website != "" {
		b.WriteString(" • " + escape(r.Website))
	}
	b.WriteString("</div></div>")
	b.WriteString("</div>")
	for _, s := range r.Sections {
		if !s.IsVisible {
			continue
		}
		b.WriteString("<div class=section>")
		b.WriteString("<div class=sec-title>" + escape(s.Title) + "</div>")
		if s.Type == "skills" {
			for _, it := range s.Items {
				if it.Description != "" {
					b.WriteString("<span class=badge>" + escape(it.Description) + "</span>")
				}
			}
		} else {
			for _, it := range s.Items {
				b.WriteString("<div class=item>")
				b.WriteString("<div class=item-top>")
				b.WriteString("<div>")
				if it.Title != "" {
					b.WriteString("<div class=item-title>" + escape(it.Title) + "</div>")
				}
				if it.Subtitle != "" || it.Location != "" {
					meta := strings.TrimSpace(strings.Join(filterEmpty([]string{it.Subtitle, it.Location}, " | "), ""))
					if meta != "" {
						b.WriteString("<div class=item-meta>" + escape(meta) + "</div>")
					}
				}
				b.WriteString("</div>")
				if it.DateRange != "" {
					b.WriteString("<div class=item-meta>" + escape(it.DateRange) + "</div>")
				}
				b.WriteString("</div>")
				if it.Description != "" {
					b.WriteString("<div class=desc>" + escape(it.Description) + "</div>")
				}
				b.WriteString("</div>")
			}
		}
		b.WriteString("</div>")
	}
	b.WriteString("</div></body></html>")
	return b.String()
}

func filterEmpty(items []string, sep string) []string {
	var out []string
	for _, v := range items {
		if strings.TrimSpace(v) != "" {
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return []string{}
	}
	return []string{strings.Join(out, sep)}
}
func escape(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(s, "<", "&lt;"), ">", "&gt;")
}

func urlEncode(s string) string {
	r := strings.NewReplacer(" ", "%20", "\n", "")
	return r.Replace(s)
}

func buildHeader(r models.Resume) string {
	return "<div style='font-size:10px;width:100%;padding:0 10mm'><span style='float:left'>" + escape(r.FullName) + "</span><span style='float:right'><span class='pageNumber'></span>/<span class='totalPages'></span></span></div>"
}

func buildFooter() string {
	return "<div style='font-size:10px;width:100%;padding:0 10mm;color:#6b7280'><span style='float:right'>Generated by OpenResume</span></div>"
}
