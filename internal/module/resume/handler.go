package resume

import (
	"encoding/json"
	"net/http"

	"openresume/internal/middleware"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	items := make([]gin.H, 0, len(list))
	for _, r := range list {
		sections := make([]gin.H, 0, len(r.Sections))
		for _, s := range r.Sections {
			previewItems := make([]gin.H, 0, len(s.Items))
			max := 3
			for i, it := range s.Items {
				if i >= max {
					break
				}
				previewItems = append(previewItems, gin.H{
					"id":          it.ExternalID,
					"title":       it.Title,
					"subtitle":    it.Subtitle,
					"dateRange":   it.DateRange,
					"description": it.Description,
				})
			}
			sections = append(sections, gin.H{
				"id":        s.ExternalID,
				"type":      s.Type,
				"title":     s.Title,
				"isVisible": s.IsVisible,
				"items":     previewItems,
			})
		}
		items = append(items, gin.H{
			"id":           r.ExternalID,
			"title":        r.Title,
			"templateId":   r.TemplateID,
			"themeConfig":  gin.H{"color": r.ThemeColor, "fontFamily": r.ThemeFont, "spacing": r.ThemeSpacing},
			"lastModified": r.LastModified,
			"personalInfo": gin.H{
				"fullName":        r.FullName,
				"jobTitle":        r.JobTitle,
				"email":           r.Email,
				"phone":           r.Phone,
				"avatarUrl":       r.AvatarURL,
				"gender":          r.Gender,
				"age":             r.Age,
				"maritalStatus":   r.MaritalStatus,
				"politicalStatus": r.PoliticalStatus,
				"birthplace":      r.Birthplace,
				"ethnicity":       r.Ethnicity,
				"height":          r.Height,
				"weight":          r.Weight,
				"customInfo": func() []gin.H {
					var items []gin.H
					if r.CustomInfo != "" {
						if arr, ok := parseCustomInfo(r.CustomInfo); ok {
							for _, it := range arr {
								items = append(items, gin.H{"label": it.Label, "value": it.Value})
							}
						}
					}
					return items
				}(),
			},
			"sections": sections,
		})
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *Handler) Create(c *gin.Context) {
	var req ResumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": res.ExternalID})
}

func (h *Handler) Get(c *gin.Context) {
	res, code, err := h.svc.GetOwnedResume(c, c.Param("id"), true)
	if err != nil {
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
	c.JSON(http.StatusOK, res)
}

func (h *Handler) Update(c *gin.Context) {
	var req ResumeReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}
	code, err := h.svc.UpdateOwnedResume(c, c.Param("id"), req)
	if err != nil {
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
