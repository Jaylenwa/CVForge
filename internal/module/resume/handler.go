package resume

import (
	"encoding/json"
	"net/http"

	"openresume/internal/middleware"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	svc *Service
}

func NewHandler() *Handler {
	return &Handler{svc: NewService()}
}

type customKV struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

func parseCustomInfo(s string) ([]customKV, bool) {
	if s == "" {
		return nil, false
	}
	var arr []customKV
	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		return nil, false
	}
	return arr, true
}

func (h *Handler) List(c *gin.Context) {
	uid, ok := middleware.UID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	list, err := h.svc.ListUserResumes(uid)
	if err != nil {
		logger.WithCtx(c).Error("resume.list failed", zap.Error(err), zap.Uint("uid", uid))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	items := make([]ResumeDTO, 0, len(list))
	for _, r := range list {
		items = append(items, ToPreviewDTO(r, 3))
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) Create(c *gin.Context) {
	var req ResumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithCtx(c).Error("resume.create bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	uid, ok := middleware.UID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	res, err := h.svc.CreateResume(uid, req)
	if err != nil {
		logger.WithCtx(c).Error("resume.create failed", zap.Error(err), zap.Uint("uid", uid))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": res.ExternalID})
}

func (h *Handler) Get(c *gin.Context) {
	res, code, err := h.svc.GetOwnedResume(c, c.Param("id"), true)
	if err != nil {
		logger.WithCtx(c).Error("resume.get failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 401:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		case 403:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.JSON(http.StatusOK, ToDTO(res))
}

func (h *Handler) Update(c *gin.Context) {
	var req ResumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.WithCtx(c).Error("resume.update bad request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	code, err := h.svc.UpdateOwnedResume(c, c.Param("id"), req)
	if err != nil {
		logger.WithCtx(c).Error("resume.update failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 401:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		case 403:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Delete(c *gin.Context) {
	code, err := h.svc.DeleteOwnedResume(c, c.Param("id"))
	if err != nil {
		logger.WithCtx(c).Error("resume.delete failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 401:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		case 403:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
