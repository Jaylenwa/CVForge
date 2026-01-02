package resume

import (
	"context"
	"encoding/json"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

type ResumeReq struct {
	Title       string `json:"title"`
	TemplateId  string `json:"templateId"`
	ThemeConfig struct {
		Color      string `json:"color"`
		FontFamily string `json:"fontFamily"`
		Spacing    string `json:"spacing"`
	} `json:"themeConfig"`
	PersonalInfo struct {
		FullName        string `json:"fullName"`
		Email           string `json:"email"`
		Phone           string `json:"phone"`
		AvatarURL       string `json:"avatarUrl"`
		JobTitle        string `json:"jobTitle"`
		Gender          string `json:"gender"`
		Age             string `json:"age"`
		MaritalStatus   string `json:"maritalStatus"`
		PoliticalStatus string `json:"politicalStatus"`
		Birthplace      string `json:"birthplace"`
		Ethnicity       string `json:"ethnicity"`
		Height          string `json:"height"`
		Weight          string `json:"weight"`
		CustomInfo      []struct {
			Label string `json:"label"`
			Value string `json:"value"`
		} `json:"customInfo"`
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
			Major       string `json:"major"`
			Degree      string `json:"degree"`
			TimeStart   string `json:"timeStart"`
			TimeEnd     string `json:"timeEnd"`
			Today       bool   `json:"today"`
			Description string `json:"description"`
		} `json:"items"`
	} `json:"sections"`
}

func (s *Service) ListUserResumes(uid uint) ([]Resume, error) {
	return s.repo.ListByUser(uid)
}

func (s *Service) CreateResume(uid uint, req ResumeReq) (Resume, error) {
	res := s.toModel(uid, req)
	err := s.repo.Create(&res)
	if err == nil && req.TemplateId != "" {
		_ = s.repo.IncrementTemplateUsage(req.TemplateId)
		_ = cache.RDB.Del(context.Background(), string(common.RedisKeyTemplatesListAll)).Err()
	}
	return res, err
}

func (s *Service) GetOwnedResume(c *gin.Context, externalID string, preload bool) (Resume, int, error) {
	uid, ok := middleware.UID(c)
	if !ok {
		return Resume{}, 401, gorm.ErrInvalidTransaction
	}
	res, err := s.repo.FindByExternal(externalID, preload)
	if err != nil {
		return Resume{}, 404, err
	}
	if res.UserID != uid {
		return Resume{}, 403, gorm.ErrInvalidData
	}
	return res, 200, nil
}

func (s *Service) UpdateOwnedResume(c *gin.Context, externalID string, req ResumeReq) (int, error) {
	uid, ok := middleware.UID(c)
	if !ok {
		return 401, gorm.ErrInvalidTransaction
	}
	existing, err := s.repo.FindByExternal(externalID, false)
	if err != nil {
		return 404, err
	}
	if existing.UserID != uid {
		return 403, gorm.ErrInvalidData
	}
	updated := s.toModel(existing.UserID, req)
	updated.ExternalID = existing.ExternalID
	updated.Model.ID = existing.Model.ID
	updated.Model.CreatedAt = existing.Model.CreatedAt
	if err := s.repo.Replace(existing, updated); err != nil {
		return 500, err
	}
	return 200, nil
}

func (s *Service) DeleteOwnedResume(c *gin.Context, externalID string) (int, error) {
	uid, ok := middleware.UID(c)
	if !ok {
		return 401, gorm.ErrInvalidTransaction
	}
	existing, err := s.repo.FindByExternal(externalID, false)
	if err != nil {
		return 404, err
	}
	if existing.UserID != uid {
		return 403, gorm.ErrInvalidData
	}
	if err := s.repo.DeleteByID(existing.ID); err != nil {
		return 500, err
	}
	return 200, nil
}

func (s *Service) toModel(uid uint, req ResumeReq) Resume {
	r := Resume{
		ExternalID:   uuid.NewString(),
		UserID:       uid,
		Title:        req.Title,
		TemplateID:   req.TemplateId,
		LastModified: time.Now().UnixMilli(),
		Personal: ResumePersonal{
			FullName:        req.PersonalInfo.FullName,
			Email:           req.PersonalInfo.Email,
			Phone:           req.PersonalInfo.Phone,
			AvatarURL:       req.PersonalInfo.AvatarURL,
			JobTitle:        req.PersonalInfo.JobTitle,
			Gender:          req.PersonalInfo.Gender,
			Age:             req.PersonalInfo.Age,
			MaritalStatus:   req.PersonalInfo.MaritalStatus,
			PoliticalStatus: req.PersonalInfo.PoliticalStatus,
			Birthplace:      req.PersonalInfo.Birthplace,
			Ethnicity:       req.PersonalInfo.Ethnicity,
			Height:          req.PersonalInfo.Height,
			Weight:          req.PersonalInfo.Weight,
			CustomInfo:      s.formatCustomInfo(req.PersonalInfo.CustomInfo),
		},
		Theme: ResumeTheme{
			Color:   req.ThemeConfig.Color,
			Font:    req.ThemeConfig.FontFamily,
			Spacing: req.ThemeConfig.Spacing,
		},
	}
	sanitize := func(s string) string {
		p := bluemonday.NewPolicy()
		p.AllowElements("p", "ul", "ol", "li", "strong", "em", "br", "span", "a")
		p.AllowAttrs("href").OnElements("a")
		p.AllowAttrs("rel").OnElements("a")
		p.AllowAttrs("target").OnElements("a")
		p.AllowStandardURLs()
		p.RequireNoFollowOnLinks(true)
		p.AddTargetBlankToFullyQualifiedLinks(true)
		return p.Sanitize(s)
	}
	for si, sct := range req.Sections {
		sec := ResumeSection{
			ExternalID: sct.Id,
			Type:       sct.Type,
			Title:      sct.Title,
			IsVisible:  sct.IsVisible,
			OrderNum:   si,
		}
		for ii, it := range sct.Items {
			sec.Items = append(sec.Items, ResumeItem{
				ExternalID:  it.Id,
				Title:       it.Title,
				Subtitle:    it.Subtitle,
				Major:       it.Major,
				Degree:      it.Degree,
				TimeStart:   it.TimeStart,
				TimeEnd:     it.TimeEnd,
				Today:       it.Today,
				Description: sanitize(it.Description),
				OrderNum:    ii,
			})
		}
		r.Sections = append(r.Sections, sec)
	}
	return r
}

func (s *Service) formatCustomInfo(items []struct {
	Label string `json:"label"`
	Value string `json:"value"`
}) string {
	if len(items) == 0 {
		return ""
	}
	b, err := json.Marshal(items)
	if err != nil {
		return ""
	}
	return string(b)
}
