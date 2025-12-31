package pdf

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"openresume/internal/infra/cache"

	"github.com/aws/aws-sdk-go-v2/aws"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

func (h *Handler) GeneratePDF(c *gin.Context) {
	pdf, code, err := h.svc.GeneratePDF(c, c.Param("id"))
	if err != nil {
		switch code {
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		case 503:
			c.JSON(http.StatusNotImplemented, gin.H{"error": "pdf service unavailable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "pdf service unavailable"})
		}
		return
	}
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=resume.pdf")
	c.Writer.Write(pdf)
}

func (h *Handler) GenerateImage(c *gin.Context) {
	img, code, err := h.svc.GenerateImage(c, c.Param("id"))
	if err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "not ready"})
		return
	}
	if token == "" || !repo.ValidateOTT(c, id, token) {
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
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=resume.pdf")
	c.Header("Cache-Control", "no-store, no-cache, must-revalidate, max-age=0, private")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")
	if strings.Contains(url, "/public/uploads/") {
		p := filepath.Join("uploads", filename)
		f, err := os.Open(p)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		defer f.Close()
		_, _ = io.Copy(c.Writer, f)
		h.svc.deleteJobFile(c, id)
		return
	}
	bucket := h.svc.sysConfig.Get("storage_s3_bucket")
	region := h.svc.sysConfig.Get("storage_s3_region")
	endpoint := h.svc.sysConfig.Get("storage_s3_endpoint")
	ak := h.svc.sysConfig.Get("storage_s3_access_key")
	sk := h.svc.sysConfig.Get("storage_s3_secret_key")
	if bucket != "" && region != "" && filename != "" {
		opts := []func(*awscfg.LoadOptions) error{awscfg.WithRegion(region)}
		if ak != "" && sk != "" {
			opts = append(opts, awscfg.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ak, sk, "")))
		}
		cfg, err := awscfg.LoadDefaultConfig(c, opts...)
		if err == nil {
			cli := s3.NewFromConfig(cfg, func(o *s3.Options) {
				if endpoint != "" {
					o.BaseEndpoint = aws.String(endpoint)
				}
			})
			obj, err := cli.GetObject(c, &s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(filename),
			})
			if err == nil {
				defer obj.Body.Close()
				_, _ = io.Copy(c.Writer, obj.Body)
			}
		}
	}
	h.svc.deleteJobFile(c, id)
}
