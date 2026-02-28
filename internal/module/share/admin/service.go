package admin

import (
	"context"
	"strconv"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

type ShareLinkItem struct {
	ID           uint
	ResumeID     uint
	UserID       uint
	UserName     string
	Slug         string
	IsPublic     bool
	HasPassword  bool
	ExpiresAt    *time.Time
	Views        uint64
	LastAccessAt *time.Time
	CreatedAt    time.Time
}

func (s *Service) AdminListShareLinks(ctx context.Context, page, size int, slug string, isPublic *bool) ([]ShareLinkItem, int64, error) {
	list, total, err := s.repo.AdminListShareLinks(page, size, slug, isPublic)
	if err != nil {
		return nil, 0, err
	}

	resumeIDs := make([]uint, 0, len(list))
	for _, sl := range list {
		if sl.ResumeID != 0 {
			resumeIDs = append(resumeIDs, sl.ResumeID)
		}
	}
	userIDByResume := make(map[uint]uint, len(resumeIDs))
	if len(resumeIDs) > 0 {
		rs, err := s.repo.FindResumesByIDs(resumeIDs)
		if err == nil {
			for _, r := range rs {
				userIDByResume[r.ID] = r.UserID
			}
		}
	}

	uidSet := make(map[uint]struct{})
	for _, sl := range list {
		uid := sl.UserID
		if uid == 0 {
			uid = userIDByResume[sl.ResumeID]
		}
		if uid != 0 {
			uidSet[uid] = struct{}{}
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

	viewsBySlug := map[string]uint64{}
	lastAccessBySlug := map[string]time.Time{}
	if cache.RDB != nil && len(list) > 0 {
		viewKeys := make([]string, 0, len(list))
		lastKeys := make([]string, 0, len(list))
		for _, s := range list {
			viewKeys = append(viewKeys, common.RedisKeyViews.F(s.Slug))
			lastKeys = append(lastKeys, common.RedisKeyShareLastAccess.F(s.Slug))
		}
		if vals, err := cache.RDB.MGet(ctx, viewKeys...).Result(); err == nil {
			for i, v := range vals {
				if v == nil || i >= len(list) {
					continue
				}
				switch t := v.(type) {
				case string:
					if n, err := strconv.ParseUint(t, 10, 64); err == nil {
						viewsBySlug[list[i].Slug] = n
					}
				}
			}
		}
		if vals, err := cache.RDB.MGet(ctx, lastKeys...).Result(); err == nil {
			for i, v := range vals {
				if v == nil || i >= len(list) {
					continue
				}
				switch t := v.(type) {
				case string:
					if ms, err := strconv.ParseInt(t, 10, 64); err == nil && ms > 0 {
						lastAccessBySlug[list[i].Slug] = time.UnixMilli(ms)
					}
				}
			}
		}
	}

	items := make([]ShareLinkItem, 0, len(list))
	for _, sl := range list {
		uid := sl.UserID
		if uid == 0 {
			uid = userIDByResume[sl.ResumeID]
		}
		views := sl.Views
		if v, ok := viewsBySlug[sl.Slug]; ok && v > 0 {
			views = v
		}
		var lastAccessAt *time.Time
		if t, ok := lastAccessBySlug[sl.Slug]; ok {
			tt := t
			lastAccessAt = &tt
		} else {
			lastAccessAt = sl.LastAccessAt
		}
		items = append(items, ShareLinkItem{
			ID:           sl.ID,
			ResumeID:     sl.ResumeID,
			UserID:       uid,
			UserName:     nameMap[uid],
			Slug:         sl.Slug,
			IsPublic:     sl.IsPublic,
			HasPassword:  sl.Password != "",
			ExpiresAt:    sl.ExpiresAt,
			Views:        views,
			LastAccessAt: lastAccessAt,
			CreatedAt:    sl.CreatedAt,
		})
	}
	return items, total, nil
}

func (s *Service) AdminUpdateShareLink(ctx context.Context, slug string, isPublic bool) error {
	sl, err := s.repo.FindShareLinkBySlug(slug)
	if err != nil {
		return err
	}
	sl.IsPublic = isPublic
	if err := s.repo.SaveShareLink(&sl); err != nil {
		return err
	}
	if cache.RDB != nil {
		_ = cache.RDB.Del(ctx, common.RedisKeyPublicResume.F(sl.Slug)).Err()
	}
	return nil
}

func (s *Service) AdminDeleteShareLink(ctx context.Context, slug string) error {
	if err := s.repo.DeleteBySlug(slug); err != nil {
		return err
	}
	if cache.RDB != nil {
		_ = cache.RDB.Del(ctx, common.RedisKeyPublicResume.F(slug)).Err()
	}
	return nil
}
