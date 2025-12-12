package handlers

import (
	"net/http"
	"time"

	"openresume/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type resumeReq struct {
	Title       string `json:"title"`
	TemplateId  string `json:"templateId"`
	ThemeConfig struct {
		Color      string `json:"color"`
		FontFamily string `json:"fontFamily"`
		Spacing    string `json:"spacing"`
	} `json:"themeConfig"`
	PersonalInfo struct {
		FullName  string `json:"fullName"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
		Website   string `json:"website"`
		AvatarURL string `json:"avatarUrl"`
		JobTitle  string `json:"jobTitle"`
	} `json:"personalInfo"`
	Sections []struct {
		Id        string `json:"id"`
		Type      string `json:"type"`
		Title     string `json:"title"`
		IsVisible bool   `json:"isVisible"`
		Items     []struct {
			Id          string `json:"id"`
			Title       string `json:"title"`
			Subtitle    string `json:"subtitle"`
			DateRange   string `json:"dateRange"`
			Location    string `json:"location"`
			Description string `json:"description"`
		} `json:"items"`
	} `json:"sections"`
}

func RegisterResumeRoutes(r *gin.RouterGroup, db *gorm.DB, auth gin.HandlerFunc) {

	r.GET("/resumes", auth, func(c *gin.Context) {
		uidVal, _ := c.Get("uid")
		var uid uint
		switch v := uidVal.(type) {
		case uint:
			uid = v
		case int:
			uid = uint(v)
		case float64:
			uid = uint(v)
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		var list []models.Resume
		if err := db.Where("user_id = ?", uid).
			Preload("Sections.Items").
			Order("updated_at desc").
			Find(&list).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		// 精简响应映射，避免过多数据
		items := make([]gin.H, 0, len(list))
		for _, r := range list {
			sections := make([]gin.H, 0, len(r.Sections))
			for _, s := range r.Sections {
				// 可按需截断预览项数量，这里保留前 3 条
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
						"location":    it.Location,
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
				"personalInfo": gin.H{"fullName": r.FullName, "jobTitle": r.JobTitle, "email": r.Email, "phone": r.Phone, "address": r.Address, "website": r.Website, "avatarUrl": r.AvatarURL},
				"sections":     sections,
			})
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	})

	r.POST("/resumes", auth, func(c *gin.Context) {
		uidVal, _ := c.Get("uid")
		var req resumeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var uid uint
		switch v := uidVal.(type) {
		case uint:
			uid = v
		case int:
			uid = uint(v)
		case float64:
			uid = uint(v)
		default:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		res := toResumeModel(uid, req)
		if err := db.Create(&res).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": res.ExternalID})
	})

	r.GET("/resumes/:id", auth, func(c *gin.Context) {
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).Preload("Sections.Items").First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		c.JSON(http.StatusOK, res)
	})

	r.PUT("/resumes/:id", auth, func(c *gin.Context) {
		var req resumeReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
			return
		}
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		updated := toResumeModel(res.UserID, req)
		// keep ExternalID
		updated.ExternalID = res.ExternalID
		updated.Model.ID = res.Model.ID
		// preserve original CreatedAt to avoid zero-value overwrite
		updated.Model.CreatedAt = res.Model.CreatedAt
		// transactional replace sections/items
		if err := db.Transaction(func(tx *gorm.DB) error {
			var secIDs []uint
			if err := tx.Model(&models.ResumeSection{}).Where("resume_id = ?", res.ID).Pluck("id", &secIDs).Error; err != nil {
				return err
			}
			if len(secIDs) > 0 {
				if err := tx.Where("section_id IN ?", secIDs).Delete(&models.ResumeItem{}).Error; err != nil {
					return err
				}
			}
			if err := tx.Where("resume_id = ?", res.ID).Delete(&models.ResumeSection{}).Error; err != nil {
				return err
			}
			if err := tx.Session(&gorm.Session{FullSaveAssociations: true}).Omit("CreatedAt").Save(&updated).Error; err != nil {
				return err
			}
			return nil
		}); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})

	r.DELETE("/resumes/:id", auth, func(c *gin.Context) {
		var res models.Resume
		if err := db.Where("external_id = ?", c.Param("id")).First(&res).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		if err := db.Delete(&res).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "db error"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true})
	})
}

func toResumeModel(uid uint, req resumeReq) models.Resume {
	r := models.Resume{
		ExternalID:   uuid.NewString(),
		UserID:       uid,
		Title:        req.Title,
		TemplateID:   req.TemplateId,
		ThemeColor:   req.ThemeConfig.Color,
		ThemeFont:    req.ThemeConfig.FontFamily,
		ThemeSpacing: req.ThemeConfig.Spacing,
		LastModified: time.Now().UnixMilli(),
		FullName:     req.PersonalInfo.FullName,
		JobTitle:     req.PersonalInfo.JobTitle,
		Email:        req.PersonalInfo.Email,
		Phone:        req.PersonalInfo.Phone,
		Address:      req.PersonalInfo.Address,
		Website:      req.PersonalInfo.Website,
		AvatarURL:    req.PersonalInfo.AvatarURL,
	}
	for si, s := range req.Sections {
		sec := models.ResumeSection{
			ExternalID: s.Id,
			Type:       s.Type,
			Title:      s.Title,
			IsVisible:  s.IsVisible,
			OrderNum:   si,
		}
		for ii, it := range s.Items {
			sec.Items = append(sec.Items, models.ResumeItem{
				ExternalID:  it.Id,
				Title:       it.Title,
				Subtitle:    it.Subtitle,
				DateRange:   it.DateRange,
				Location:    it.Location,
				Description: it.Description,
				OrderNum:    ii,
			})
		}
		r.Sections = append(r.Sections, sec)
	}
	return r
}
