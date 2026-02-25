package resume

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

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

func (h *Handler) List(c *gin.Context) {
	uid, ok := middleware.UID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	list, err := h.svc.ListUserResumes(uid)
	if err != nil {
		logger.WithCtx(c).Error("resume.list failed", zap.Error(err), zap.Uint("uid", uid))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		var inv *InvalidSectionTypeError
		if errors.As(err, &inv) {
			c.JSON(http.StatusBadRequest, gin.H{"error": inv.Error(), "type": inv.Value})
			return
		}
		logger.WithCtx(c).Error("resume.create failed", zap.Error(err), zap.Uint("uid", uid))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": res.ID})
}

func (h *Handler) Get(c *gin.Context) {
	id, parseErr := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	res, code, err := h.svc.GetOwnedResume(c, uint(id), true)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	id, parseErr := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	code, err := h.svc.UpdateOwnedResume(c, uint(id), req)
	if err != nil {
		logger.WithCtx(c).Error("resume.update failed", zap.Error(err), zap.Int("code", code), zap.String("id", c.Param("id")))
		switch code {
		case 400:
			var inv *InvalidSectionTypeError
			if errors.As(err, &inv) {
				c.JSON(http.StatusBadRequest, gin.H{"error": inv.Error(), "type": inv.Value})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		case 401:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		case 403:
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		case 404:
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *Handler) Delete(c *gin.Context) {
	id, parseErr := strconv.ParseUint(strings.TrimSpace(c.Param("id")), 10, 64)
	if parseErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	code, err := h.svc.DeleteOwnedResume(c, uint(id))
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
