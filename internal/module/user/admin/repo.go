package admin

import (
	"errors"
	"strings"

	"cvforge/internal/common"
	"cvforge/internal/infra/database"
	usermod "cvforge/internal/module/user"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func DefaultRepo() *Repo {
	return &Repo{db: database.DB}
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{db: db}
}

type UserProviderRow struct {
	UserID   uint                `gorm:"column:user_id"`
	Provider common.ProviderType `gorm:"column:provider"`
}

func (r *Repo) ListUsers(page, size int, emailQ, nameQ, role string, isActive *bool) ([]usermod.User, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	var list []usermod.User
	q := r.db.Model(&usermod.User{})
	if emailQ = strings.TrimSpace(emailQ); emailQ != "" || strings.TrimSpace(nameQ) != "" {
		nameQ = strings.TrimSpace(nameQ)
		if emailQ != "" && nameQ != "" {
			q = q.Where("(email LIKE ? OR name LIKE ?)", "%"+emailQ+"%", "%"+nameQ+"%")
		} else if emailQ != "" {
			q = q.Where("email LIKE ?", "%"+emailQ+"%")
		} else {
			q = q.Where("name LIKE ?", "%"+nameQ+"%")
		}
	}
	if role = strings.TrimSpace(role); role != "" {
		q = q.Where("role = ?", role)
	}
	if isActive != nil {
		q = q.Where("is_active = ?", *isActive)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Order("id desc").Offset((page - 1) * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *Repo) ListOAuthProviders(userIDs []uint) ([]UserProviderRow, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	var rows []UserProviderRow
	err := r.db.Model(&usermod.OAuthAccount{}).Select("user_id, provider").Where("user_id IN ?", userIDs).Find(&rows).Error
	return rows, err
}

func (r *Repo) UpdateUserActive(id uint, active bool) error {
	return r.db.Model(&usermod.User{}).Where("id = ?", id).Update("is_active", active).Error
}

func (r *Repo) UpdateUserPasswordHash(id uint, hash string) error {
	return r.db.Model(&usermod.User{}).Where("id = ?", id).Update("password_hash", hash).Error
}

func (r *Repo) CreateAuditLog(al *usermod.AuditLog) error {
	if al == nil {
		return errors.New("invalid")
	}
	return r.db.Create(al).Error
}

