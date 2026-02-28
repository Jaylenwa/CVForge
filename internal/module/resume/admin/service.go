package admin

import (
	"context"
	"errors"
	"strconv"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	resumemod "cvforge/internal/module/resume"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service struct {
	repo     *Repo
	userRepo *resumemod.Repo
}

func NewService() *Service {
	return &Service{
		repo:     DefaultRepo(),
		userRepo: resumemod.DefaultRepo(),
	}
}

type AuditActor struct {
	ActorID uint
	IP      string
	UA      string
}

type ResumeListItem struct {
	Resume    resumemod.ResumeDTO `json:"resume"`
	UserID    uint               `json:"userId"`
	UserName  string             `json:"userName"`
	CreatedAt any                `json:"createdAt"`
	UpdatedAt any                `json:"updatedAt"`
}

type ResumeListResult struct {
	Items      []ResumeListItem `json:"items"`
	Page       int              `json:"page"`
	PageSize   int              `json:"pageSize"`
	Total      int64            `json:"total"`
	TotalPages int              `json:"totalPages"`
	HasNext    bool             `json:"hasNext"`
}

func (s *Service) AdminListResumes(page, size int, userID *uint, title, templateID string) (ResumeListResult, error) {
	list, total, err := s.repo.ListResumes(page, size, userID, title, templateID)
	if err != nil {
		return ResumeListResult{}, err
	}
	uidSet := make(map[uint]struct{})
	for _, r := range list {
		if r.UserID != 0 {
			uidSet[r.UserID] = struct{}{}
		}
	}
	uids := make([]uint, 0, len(uidSet))
	for id := range uidSet {
		uids = append(uids, id)
	}
	nameMap := make(map[uint]string, len(uids))
	if len(uids) > 0 {
		users, err := s.repo.FindUsersByIDs(uids)
		if err == nil {
			for _, u := range users {
				if u.Name != "" {
					nameMap[u.ID] = u.Name
				} else if u.Email != nil {
					nameMap[u.ID] = *u.Email
				}
			}
		}
	}
	items := make([]ResumeListItem, 0, len(list))
	for _, r := range list {
		items = append(items, ResumeListItem{
			Resume:    resumemod.ToDTO(r),
			UserID:    r.UserID,
			UserName:  nameMap[r.UserID],
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	totalPages := (int(total) + size - 1) / size
	hasNext := page*size < int(total)
	return ResumeListResult{
		Items:      items,
		Page:       page,
		PageSize:   size,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    hasNext,
	}, nil
}

func (s *Service) AdminGetResume(id uint) (resumemod.ResumeDTO, error) {
	res, err := s.repo.FindResumeFullByID(id)
	if err != nil {
		return resumemod.ResumeDTO{}, err
	}
	return resumemod.ToDTO(res), nil
}

func (s *Service) AdminDeleteResume(actor AuditActor, id uint) error {
	if id == 0 {
		return errors.New("invalid")
	}
	if err := s.repo.DeleteResumeByID(id); err != nil {
		return err
	}
	_ = s.repo.CreateAuditLog(&resumemod.AuditLog{
		ActorID:    actor.ActorID,
		Action:     "resume.delete",
		TargetType: "resume",
		TargetID:   strconv.FormatUint(uint64(id), 10),
		Metadata:   "",
		IP:         actor.IP,
		UA:         actor.UA,
	})
	return nil
}

func (s *Service) AdminUpdateResumeVisibility(ctx context.Context, actor AuditActor, id uint, isPublic bool) (string, error) {
	if id == 0 {
		return "", errors.New("invalid")
	}
	res, err := s.repo.FindResumeByID(id)
	if err != nil {
		return "", err
	}
	sl, err := s.repo.FindShareLinkByResumeID(res.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			sl = resumemod.ShareLink{ResumeID: res.ID, Slug: uuid.NewString()[:8], IsPublic: isPublic}
			if err := s.repo.CreateShareLink(&sl); err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		sl.IsPublic = isPublic
		if err := s.repo.SaveShareLink(&sl); err != nil {
			return "", err
		}
	}
	if cache.RDB != nil {
		_ = cache.RDB.Del(ctx, common.RedisKeyPublicResume.F(sl.Slug)).Err()
	}
	_ = s.repo.CreateAuditLog(&resumemod.AuditLog{
		ActorID:    actor.ActorID,
		Action:     "resume.visibility",
		TargetType: "resume",
		TargetID:   strconv.FormatUint(uint64(id), 10),
		Metadata:   strconv.FormatBool(isPublic),
		IP:         actor.IP,
		UA:         actor.UA,
	})
	return sl.Slug, nil
}

