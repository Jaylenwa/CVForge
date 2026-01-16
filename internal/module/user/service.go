package user

import (
	"context"
	"errors"
	conf "openresume/internal/module/config"
	"openresume/internal/pkg/storage"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repo
}

func NewService() *Service {
	return &Service{repo: DefaultRepo()}
}

func (s *Service) GetMe(id any) (User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) UpdateProfile(id any, name, avatar, lang string) error {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	oldAvatar := u.AvatarURL
	if name != "" {
		u.Name = name
	}
	if avatar != "" {
		u.AvatarURL = avatar
	}
	if lang != "" {
		u.Language = lang
	}
	if err := s.repo.Save(&u); err != nil {
		return err
	}
	if avatar != "" && oldAvatar != "" && avatar != oldAvatar {
		sys := conf.NewService()
		cfg := sys.GetStorageSettings()
		if up, e := storage.NewFromSettings(cfg); e == nil {
			_ = up.Delete(context.Background(), oldAvatar)
		}
	}
	return nil
}

func (s *Service) UpdatePassword(id any, current, newpw string) error {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if newpw == "" {
		return errors.New("invalid")
	}
	if u.PasswordHash != "" {
		if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(current)); err != nil {
			return gorm.ErrInvalidData
		}
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newpw), bcrypt.DefaultCost)
	u.PasswordHash = string(hash)
	return s.repo.Save(&u)
}
