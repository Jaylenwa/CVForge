package admin

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	usermod "cvforge/internal/module/user"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo     *Repo
	userRepo *usermod.Repo
}

func NewService() *Service {
	return &Service{
		repo:     DefaultRepo(),
		userRepo: usermod.DefaultRepo(),
	}
}

type AuditActor struct {
	ActorID uint
	IP      string
	UA      string
}

func (s *Service) writeAudit(actor AuditActor, action, targetType, targetID, metadata string) {
	_ = s.repo.CreateAuditLog(&usermod.AuditLog{
		ActorID:    actor.ActorID,
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Metadata:   metadata,
		IP:         actor.IP,
		UA:         actor.UA,
	})
}

type AdminUserItem struct {
	ID            uint        `json:"id"`
	Email         string      `json:"email"`
	Name          string      `json:"name"`
	AvatarURL     string      `json:"avatarUrl"`
	Language      string      `json:"language"`
	Role          common.Role `json:"role"`
	IsActive      bool        `json:"isActive"`
	LastLoginAt   *time.Time  `json:"lastLoginAt"`
	CreatedAt     time.Time   `json:"createdAt"`
	UpdatedAt     time.Time   `json:"updatedAt,omitempty"`
	Providers     []string    `json:"providers"`
	LoginProvider string      `json:"loginProvider"`
}

func orderedProviders(rows []UserProviderRow) map[uint][]common.ProviderType {
	sets := make(map[uint]map[common.ProviderType]struct{})
	for _, r := range rows {
		if r.UserID == 0 || r.Provider == "" {
			continue
		}
		s := sets[r.UserID]
		if s == nil {
			s = map[common.ProviderType]struct{}{}
			sets[r.UserID] = s
		}
		s[r.Provider] = struct{}{}
	}

	out := make(map[uint][]common.ProviderType, len(sets))
	for uid, s := range sets {
		ordered := make([]common.ProviderType, 0, len(s))
		if _, ok := s[common.WechatMP]; ok {
			ordered = append(ordered, common.WechatMP)
			delete(s, common.WechatMP)
		}
		if _, ok := s[common.Github]; ok {
			ordered = append(ordered, common.Github)
			delete(s, common.Github)
		}
		for p := range s {
			ordered = append(ordered, p)
		}
		out[uid] = ordered
	}
	return out
}

func (s *Service) AdminListUsers(page, size int, emailQ, nameQ, role string, isActive *bool) ([]AdminUserItem, int64, error) {
	list, total, err := s.repo.ListUsers(page, size, emailQ, nameQ, role, isActive)
	if err != nil {
		return nil, 0, err
	}
	userIDs := make([]uint, 0, len(list))
	for _, u := range list {
		userIDs = append(userIDs, u.ID)
	}
	rows, err := s.repo.ListOAuthProviders(userIDs)
	if err != nil {
		rows = nil
	}
	providersMap := orderedProviders(rows)

	items := make([]AdminUserItem, 0, len(list))
	for _, u := range list {
		email := ""
		if u.Email != nil {
			email = *u.Email
		}
		providers := providersMap[u.ID]
		providerStrs := make([]string, 0, len(providers))
		for _, p := range providers {
			providerStrs = append(providerStrs, string(p))
		}
		loginProvider := ""
		if strings.TrimSpace(email) != "" {
			loginProvider = "email"
		} else if len(providers) > 0 {
			loginProvider = string(providers[0])
		}
		items = append(items, AdminUserItem{
			ID:            u.ID,
			Email:         email,
			Name:          u.Name,
			AvatarURL:     u.AvatarURL,
			Language:      u.Language,
			Role:          u.Role,
			IsActive:      u.IsActive,
			LastLoginAt:   u.LastLoginAt,
			CreatedAt:     u.CreatedAt,
			Providers:     providerStrs,
			LoginProvider: loginProvider,
		})
	}
	return items, total, nil
}

func (s *Service) AdminGetUser(id uint) (AdminUserItem, error) {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		return AdminUserItem{}, err
	}
	email := ""
	if u.Email != nil {
		email = *u.Email
	}
	rows, err := s.repo.ListOAuthProviders([]uint{u.ID})
	if err != nil {
		rows = nil
	}
	providersMap := orderedProviders(rows)
	providers := providersMap[u.ID]
	providerStrs := make([]string, 0, len(providers))
	for _, p := range providers {
		providerStrs = append(providerStrs, string(p))
	}
	loginProvider := ""
	if strings.TrimSpace(email) != "" {
		loginProvider = "email"
	} else if len(providers) > 0 {
		loginProvider = string(providers[0])
	}
	return AdminUserItem{
		ID:            u.ID,
		Email:         email,
		Name:          u.Name,
		AvatarURL:     u.AvatarURL,
		Language:      u.Language,
		Role:          u.Role,
		IsActive:      u.IsActive,
		LastLoginAt:   u.LastLoginAt,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Providers:     providerStrs,
		LoginProvider: loginProvider,
	}, nil
}

type AdminPatchInput struct {
	Name      *string
	AvatarURL *string
	Language  *string
	Role      *common.Role
	IsActive  *bool
}

func (s *Service) AdminPatchUser(actor AuditActor, id uint, in AdminPatchInput) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if in.Name != nil {
		u.Name = *in.Name
	}
	if in.AvatarURL != nil {
		u.AvatarURL = *in.AvatarURL
	}
	if in.Language != nil {
		u.Language = *in.Language
	}
	if in.Role != nil {
		r := *in.Role
		if r == common.RoleUser || r == common.RoleAdmin {
			u.Role = r
		} else {
			return errors.New("invalid role")
		}
	}
	if in.IsActive != nil {
		u.IsActive = *in.IsActive
	}
	if err := s.userRepo.Save(&u); err != nil {
		return err
	}
	s.writeAudit(actor, "user.update", "user", strconv.FormatUint(uint64(u.ID), 10), "")
	return nil
}

func (s *Service) AdminResetPassword(actor AuditActor, id uint, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" {
		return errors.New("invalid")
	}
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err := s.repo.UpdateUserPasswordHash(u.ID, string(hash)); err != nil {
		return err
	}
	s.writeAudit(actor, "user.reset_password", "user", strconv.FormatUint(uint64(u.ID), 10), "")
	return nil
}

func (s *Service) AdminSetActive(actor AuditActor, id uint, active bool) error {
	if err := s.repo.UpdateUserActive(id, active); err != nil {
		return err
	}
	if cache.RDB != nil {
		val := "0"
		if active {
			val = "1"
		}
		_ = cache.RDB.Set(context.Background(), common.RedisKeyUserActive.F(strconv.FormatUint(uint64(id), 10)), val, 10*time.Minute).Err()
	}
	if active {
		s.writeAudit(actor, "user.unban", "user", strconv.FormatUint(uint64(id), 10), "")
	} else {
		s.writeAudit(actor, "user.ban", "user", strconv.FormatUint(uint64(id), 10), "")
	}
	return nil
}

