package resume

import (
	"context"
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
	Title      string            `json:"Title"`
	TemplateID string            `json:"TemplateID"`
	Personal   ResumePersonalDTO `json:"Personal"`
	Job        ResumeJobDTO      `json:"Job"`
	Theme      ResumeThemeDTO    `json:"Theme"`
	Sections   []SectionDTO      `json:"Sections"`
}

func (s *Service) ListUserResumes(uid uint) ([]Resume, error) {
	return s.repo.ListByUser(uid)
}

func (s *Service) CreateResume(uid uint, req ResumeReq) (Resume, error) {
	res := s.toModel(uid, req)
	err := s.repo.Create(&res)
	if err == nil && req.TemplateID != "" {
		_ = s.repo.IncrementTemplateUsage(req.TemplateID)
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
		TemplateID:   req.TemplateID,
		LastModified: time.Now().UnixMilli(),
		Personal: ResumePersonal{
			FullName:        req.Personal.FullName,
			Email:           req.Personal.Email,
			Phone:           req.Personal.Phone,
			AvatarURL:       req.Personal.AvatarURL,
			JobTitle:        req.Personal.JobTitle,
			Gender:          req.Personal.Gender,
			Age:             req.Personal.Age,
			MaritalStatus:   req.Personal.MaritalStatus,
			PoliticalStatus: req.Personal.PoliticalStatus,
			Birthplace:      req.Personal.Birthplace,
			Ethnicity:       req.Personal.Ethnicity,
			Height:          req.Personal.Height,
			Weight:          req.Personal.Weight,
			CustomInfo:      req.Personal.CustomInfo,
		},
		Job: ResumeJob{
			Job:      req.Job.Job,
			City:     req.Job.City,
			Money:    req.Job.Money,
			JoinTime: req.Job.JoinTime,
		},
		Theme: ResumeTheme{
			Color:   req.Theme.Color,
			Font:    req.Theme.Font,
			Spacing: req.Theme.Spacing,
			FontSize: req.Theme.FontSize,
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
			ExternalID: sct.ExternalID,
			Type:       sct.Type,
			Title:      sct.Title,
			IsVisible:  sct.IsVisible,
			OrderNum:   si,
		}
		for ii, it := range sct.Items {
			sec.Items = append(sec.Items, ResumeItem{
				ExternalID:  it.ExternalID,
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
