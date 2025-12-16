package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"openresume/config"
	"openresume/storage"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RegisterUploadRoutes(r *gin.RouterGroup, auth gin.HandlerFunc) {
	r.POST("/upload/avatar", auth, func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "no file"})
			return
		}
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type"})
			return
		}
		if file.Size > 2*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
			return
		}
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "open error"})
			return
		}
		defer f.Close()
		b, err := io.ReadAll(f)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "read error"})
			return
		}
		mt := mimetype.Detect(b)
		if mt == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid content"})
			return
		}
		mts := mt.String()
		if mts != "image/jpeg" && mts != "image/png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid mime"})
			return
		}
		name := uuid.NewString() + ext
		up, e := storage.New(config.Load())
		if e != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload init error"})
			return
		}
		url, err := up.Upload(c, name, b)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "upload error"})
			return
		}
		proto := c.GetHeader("X-Forwarded-Proto")
		if proto == "" {
			proto = "http"
		}
		host := c.GetHeader("X-Forwarded-Host")
		if host == "" {
			host = c.Request.Host
		}
		abs := url
		if strings.HasPrefix(url, "/") {
			abs = proto + "://" + host + url
		}
		c.JSON(http.StatusOK, gin.H{"url": abs})
	})
}
