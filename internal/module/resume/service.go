package resume

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/middleware"
	"openresume/internal/module/library"
	"openresume/internal/module/preset"

	"github.com/gin-gonic/gin"
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
	VariantID  *uint             `json:"VariantID"`
	PresetID   *uint             `json:"PresetID"`
	RoleID     *uint             `json:"RoleID"`
	Language   string            `json:"Language"`
	Personal   ResumePersonalDTO `json:"Personal"`
	Theme      ResumeThemeDTO    `json:"Theme"`
	Sections   []SectionReq      `json:"Sections"`
}

type SectionReq struct {
	Type      string    `json:"Type"`
	Title     string    `json:"Title"`
	IsVisible bool      `json:"IsVisible"`
	OrderNum  int       `json:"OrderNum"`
	Items     []ItemReq `json:"Items"`
}

type ItemReq struct {
	Title       string `json:"Title"`
	Subtitle    string `json:"Subtitle"`
	Major       string `json:"Major"`
	Degree      string `json:"Degree"`
	TimeStart   string `json:"TimeStart"`
	TimeEnd     string `json:"TimeEnd"`
	Today       bool   `json:"Today"`
	Description string `json:"Description"`
	OrderNum    int    `json:"OrderNum"`
}

func (s *Service) ListUserResumes(uid uint) ([]Resume, error) {
	return s.repo.ListByUser(uid)
}

func (s *Service) CreateResume(uid uint, req ResumeReq) (Resume, error) {
	res := s.toModel(uid, req)
	err := s.repo.Create(&res)
	if err == nil && req.VariantID != nil && *req.VariantID != 0 {
		_ = library.DefaultRepo().IncrementTemplateVariantUsage(*req.VariantID)
	} else if err == nil && req.TemplateID != "" {
		_ = s.repo.IncrementTemplateUsage(req.TemplateID)
		_ = cache.RDB.Del(context.Background(), string(common.RedisKeyTemplatesListAll)).Err()
	}
	return res, err
}

type CreateFromVariantReq struct {
	VariantID uint   `json:"VariantID"`
	Language  string `json:"Language"`
	Title     string `json:"Title"`
}

func (s *Service) CreateResumeFromVariant(uid uint, req CreateFromVariantReq) (Resume, error) {
	if req.VariantID == 0 {
		return Resume{}, errors.New("variant required")
	}
	libRepo := library.DefaultRepo()
	variant, err := libRepo.GetTemplateVariantByID(req.VariantID)
	if err != nil {
		return Resume{}, err
	}
	pRepo := preset.DefaultRepo()
	presetItem, err := pRepo.GetByIDActive(variant.PresetID)
	if err != nil {
		return Resume{}, err
	}

	type presetPayload struct {
		Title    string            `json:"title"`
		Language string            `json:"language"`
		Personal ResumePersonalDTO `json:"Personal"`
		Theme    ResumeThemeDTO    `json:"Theme"`
		Sections []struct {
			ID        string `json:"id"`
			Type      string `json:"type"`
			Title     string `json:"title"`
			IsVisible bool   `json:"isVisible"`
			Items     []struct {
				ID          string `json:"id"`
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

	payload := presetPayload{}
	if presetItem.DataJSON != "" {
		_ = json.Unmarshal([]byte(presetItem.DataJSON), &payload)
	}

	lang := req.Language
	if lang == "" {
		lang = payload.Language
	}
	if lang == "" {
		lang = "zh"
	}

	title := req.Title
	if title == "" {
		title = payload.Title
	}
	if title == "" {
		title = variant.Name
	}

	out := ResumeReq{
		Title:      title,
		TemplateID: variant.LayoutTemplateExternalID,
		VariantID:  &variant.ID,
		PresetID:   &presetItem.ID,
		RoleID:     &variant.RoleID,
		Language:   lang,
		Personal:   payload.Personal,
		Theme:      payload.Theme,
	}
	for _, sct := range payload.Sections {
		sec := SectionReq{
			Type:      sct.Type,
			Title:     sct.Title,
			IsVisible: sct.IsVisible,
		}
		for _, it := range sct.Items {
			sec.Items = append(sec.Items, ItemReq{
				Title:       it.Title,
				Subtitle:    it.Subtitle,
				Major:       it.Major,
				Degree:      it.Degree,
				TimeStart:   it.TimeStart,
				TimeEnd:     it.TimeEnd,
				Today:       it.Today,
				Description: it.Description,
			})
		}
		out.Sections = append(out.Sections, sec)
	}

	res, err := s.CreateResume(uid, out)
	return res, err
}

func (s *Service) GetOwnedResume(c *gin.Context, id uint, preload bool) (Resume, int, error) {
	uid, ok := middleware.UID(c)
	if !ok {
		return Resume{}, 401, gorm.ErrInvalidTransaction
	}
	res, err := s.repo.FindByID(id, preload)
	if err != nil {
		return Resume{}, 404, err
	}
	if res.UserID != uid {
		return Resume{}, 403, gorm.ErrInvalidData
	}
	return res, 200, nil
}

func (s *Service) UpdateOwnedResume(c *gin.Context, id uint, req ResumeReq) (int, error) {
	uid, ok := middleware.UID(c)
	if !ok {
		return 401, gorm.ErrInvalidTransaction
	}
	existing, err := s.repo.FindByID(id, false)
	if err != nil {
		return 404, err
	}
	if existing.UserID != uid {
		return 403, gorm.ErrInvalidData
	}
	updated := s.toModel(existing.UserID, req)
	updated.Model.ID = existing.Model.ID
	updated.Model.CreatedAt = existing.Model.CreatedAt
	if err := s.repo.Replace(existing, updated); err != nil {
		return 500, err
	}
	return 200, nil
}

func (s *Service) DeleteOwnedResume(c *gin.Context, id uint) (int, error) {
	uid, ok := middleware.UID(c)
	if !ok {
		return 401, gorm.ErrInvalidTransaction
	}
	existing, err := s.repo.FindByID(id, false)
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
		UserID:       uid,
		Title:        req.Title,
		TemplateID:   req.TemplateID,
		VariantID:    req.VariantID,
		PresetID:     req.PresetID,
		RoleID:       req.RoleID,
		Language:     req.Language,
		LastModified: time.Now().UnixMilli(),
		Personal: ResumePersonal{
			FullName:   req.Personal.FullName,
			Email:      req.Personal.Email,
			Phone:      req.Personal.Phone,
			AvatarURL:  req.Personal.AvatarURL,
			Job:        req.Personal.Job,
			City:       req.Personal.City,
			Money:      req.Personal.Money,
			JoinTime:   req.Personal.JoinTime,
			Gender:     req.Personal.Gender,
			Age:        req.Personal.Age,
			Degree:     req.Personal.Degree,
			CustomInfo: req.Personal.CustomInfo,
		},
		Theme: ResumeTheme{
			Color:    req.Theme.Color,
			Font:     req.Theme.Font,
			Spacing:  req.Theme.Spacing,
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
			Type:       sct.Type,
			Title:      sct.Title,
			IsVisible:  sct.IsVisible,
			OrderNum:   si,
		}
		for ii, it := range sct.Items {
			sec.Items = append(sec.Items, ResumeItem{
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
