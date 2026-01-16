package pdf

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"openresume/internal/infra/cache"
	"openresume/internal/infra/database"
	"openresume/internal/models"
	"openresume/internal/pkg/logger"
	"openresume/internal/pkg/storage"

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
	img, code, err := h.svc.GenerateImage(c, c.Param("id"))
	if err != nil {
		logger.WithCtx(c).Error("image.generate failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case 503:
			c.JSON(http.StatusNotImplemented, gin.H{"error": "image service unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "image service unavailable"})
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
	resumeID := c.Query("resumeId")
	if resumeID == "" {
		resumeID = c.Param("id")
	}
	if resumeID == "" {
		logger.WithCtx(c).Error("pdf.submit_export bad request")
		c.JSON(http.StatusBadRequest, gin.H{"error": "resumeId required"})
		return
	}
	authHeader := c.GetHeader("Authorization")
	job := ExportJob{
		ID:       uuid.NewString(),
		UserID:   userID,
		ResumeID: resumeID,
		Token:    strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer ")),
	}
	repo := NewExportRepo(h.svc)
	if err := repo.Enqueue(c, job); err != nil {
		logger.WithCtx(c).Error("pdf.submit_export enqueue failed", zap.Error(err))
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "queue unavailable"})
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
	filename := cache.RDB.Get(c, jobKey(id)+":filename").Val()
	if filename == "" {
		if idx := strings.LastIndex(url, "/"); idx >= 0 && idx+1 < len(url) {
			filename = url[idx+1:]
		}
	}
	job, _ := repo.GetJob(c, id)
	var r models.Resume
	if job.ResumeID != "" {
		_ = database.DB.Where("external_id = ?", job.ResumeID).First(&r).Error
	}
	c.Header("Content-Type", "application/pdf")
	disp := sanitizeFilename(r.Title)
	if disp == "" {
		disp = "resume"
	}
	c.Header("Content-Disposition", "attachment; filename="+disp)
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	cfg := h.svc.sysConfig.GetStorageSettings()
	up, err := storage.NewFromSettings(cfg)
	if err != nil {
		logger.WithCtx(c).Error("pdf.export_download storage init failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "storage unavailable"})
		return
	}
	rc, err := up.Download(c, url)
	if err != nil {
		logger.WithCtx(c).Error("pdf.export_download download failed", zap.Error(err), zap.String("url", url))
		c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
		return
	}
	defer rc.Close()
	_, _ = io.Copy(c.Writer, rc)
	h.svc.deleteJobFile(c, id)
}
