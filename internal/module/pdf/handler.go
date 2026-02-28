package pdf

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"cvforge/internal/infra/cache"
	"cvforge/internal/pkg/logger"
	"cvforge/internal/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

func sanitizeFilename(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "resume"
	}
	out := make([]rune, 0, len(s))
	for _, r := range s {
		if r == ' ' || r == '-' || r == '_' || r == '.' || r == '(' || r == ')' || r == '（' || r == '）' {
			out = append(out, r)
			continue
		}
		if r >= '0' && r <= '9' {
			out = append(out, r)
			continue
		}
		if r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' {
			out = append(out, r)
			continue
		}
		if r >= 0x4E00 && r <= 0x9FFF {
			out = append(out, r)
			continue
		}
	}
	if len(out) == 0 {
		return "resume"
	}
	return strings.TrimSpace(string(out))
}

func (h *Handler) GenerateImage(c *gin.Context) {
	id, parseErr := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	authHeader := c.GetHeader("Authorization")
	token := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
	img, code, err := h.svc.GenerateImageWithToken(uint(id), token)
	if err != nil {
		logger.WithCtx(c).Error("image.generate failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case 503:
			c.JSON(http.StatusNotImplemented, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", "attachment; filename=resume.png")
	c.Writer.Write(img)
}

func (h *Handler) SubmitExport(c *gin.Context) {
	uidVal, _ := c.Get("uid")
	userID := ""
	if uid, ok := uidVal.(uint); ok {
		userID = fmt.Sprintf("%d", uid)
	}
	resumeIDStr := strings.TrimSpace(c.Query("resumeId"))
	if resumeIDStr == "" {
		resumeIDStr = strings.TrimSpace(c.Param("id"))
	}
	id64, err := strconv.ParseUint(resumeIDStr, 10, 64)
	if resumeIDStr == "" || err != nil {
		logger.WithCtx(c).Error("pdf.submit_export bad request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "resumeId required"})
		return
	}
	authHeader := c.GetHeader("Authorization")
	job := ExportJob{
		ID:       uuid.NewString(),
		UserID:   userID,
		ResumeID: uint(id64),
		Token:    strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer ")),
	}
	repo := NewExportRepo(h.svc)
	if err := repo.Enqueue(c, job); err != nil {
		logger.WithCtx(c).Error("pdf.submit_export enqueue failed", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"job_id": job.ID})
}

func (h *Handler) ExportStatus(c *gin.Context) {
	id := c.Param("job_id")
	repo := NewExportRepo(h.svc)
	st, _, errMsg, err := repo.GetStatus(c, id)
	if err != nil {
		logger.WithCtx(c).Error("pdf.export_status not found", zap.Error(err), zap.String("id", id))
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if st == ExportStatusDone {
		tok, _ := repo.EnsureOTT(c, id)
		c.JSON(http.StatusOK, gin.H{"status": st, "token": tok, "error": errMsg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": st, "error": errMsg})
}

func (h *Handler) ExportDownload(c *gin.Context) {
	id := c.Param("job_id")
	token := c.Query("token")
	repo := NewExportRepo(h.svc)
	st, url, _, err := repo.GetStatus(c, id)
	if err != nil || st != ExportStatusDone || url == "" {
		logger.WithCtx(c).Error("pdf.export_download not ready", zap.Error(err), zap.String("id", id), zap.String("url", url))
		c.JSON(http.StatusBadRequest, gin.H{"error": "not ready"})
		return
	}
	if token == "" || !repo.ValidateOTT(c, id, token) {
		logger.WithCtx(c).Error("pdf.export_download invalid token", zap.String("id", id))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}
	key := jobKey(id) + ":consumed"
	if ok := h.svc.tryMarkConsumed(c, key); !ok {
		c.JSON(http.StatusGone, gin.H{"error": "already downloaded"})
		return
	}
	jobFilename := cache.RDB.Get(c, jobKey(id)+":filename").Val()
	if jobFilename == "" {
		if idx := strings.LastIndex(url, "/"); idx >= 0 && idx+1 < len(url) {
			jobFilename = url[idx+1:]
		}
	}
	job, _ := repo.GetJob(c, id)
	title := ""
	if job.ResumeID != 0 {
		if t, err := h.svc.ResumeTitle(job.ResumeID); err == nil {
			title = t
		}
	}
	c.Header("Content-Type", "application/pdf")
	disp := sanitizeFilename(title)
	if disp == "" {
		disp = sanitizeFilename(jobFilename)
	}
	if disp == "" {
		disp = "resume"
	}
	if !strings.HasSuffix(strings.ToLower(disp), ".pdf") {
		disp += ".pdf"
	}
	c.Header("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": disp}))
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	cfg := h.svc.sysConfig.GetStorageSettings()
	up, err := storage.NewFromSettings(cfg)
	if err != nil {
		logger.WithCtx(c).Error("pdf.export_download storage init failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rc, err := up.Download(c.Request.Context(), url)
	if err != nil {
		logger.WithCtx(c).Error("pdf.export_download download failed", zap.Error(err), zap.String("url", url))
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer rc.Close()
	_, _ = io.Copy(c.Writer, rc)
	h.svc.deleteJobFile(c, id)
}
