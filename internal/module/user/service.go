package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo *Repo
}

func NewService(db *gorm.DB) *Service {
	return &Service{repo: NewRepo(db)}
}

func (s *Service) GetMe(id any) (User, error) {
	return s.repo.FindByID(id)
}

func (s *Service) UpdateProfile(id any, name, avatar, lang string) error {
	u, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if name != "" {
		u.Name = name
	}
	if avatar != "" {
		u.AvatarURL = avatar
	}
	if lang != "" {
		u.Language = lang
	}
	return s.repo.Save(&u)
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
