package auth

import (
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/database"

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

func (r *Repo) GetUserTokenVersion(uid uint) (int, error) {
	var u User
	if err := r.db.Select("id, token_version").First(&u, uid).Error; err != nil {
		return 0, err
	}
	if u.TokenVersion <= 0 {
		return 1, nil
	}
	return u.TokenVersion, nil
}

func (r *Repo) SetUserLastLoginAt(uid uint, t time.Time) error {
	if uid == 0 {
		return nil
	}
	return r.db.Model(&User{}).Where("id = ?", uid).Update("last_login_at", &t).Error
}

func (r *Repo) CountUsers() (int64, error) {
	var n int64
	err := r.db.Model(&User{}).Count(&n).Error
	return n, err
}

func (r *Repo) CreateUser(u *User) error {
	return r.db.Create(u).Error
}

func (r *Repo) FindUserByEmail(email string) (User, error) {
	var u User
	err := r.db.Where("email = ?", email).First(&u).Error
	return u, err
}

func (r *Repo) FindUserByID(uid uint) (User, error) {
	var u User
	err := r.db.First(&u, uid).Error
	return u, err
}

func (r *Repo) IncrementUserTokenVersion(uid uint) error {
	if uid == 0 {
		return nil
	}
	return r.db.Model(&User{}).Where("id = ?", uid).Update("token_version", gorm.Expr("token_version + 1")).Error
}

func (r *Repo) FindOAuthAccountByProviderUnionID(provider common.ProviderType, unionID string) (OAuthAccount, error) {
	var oa OAuthAccount
	err := r.db.Where("provider_union_id <> '' AND provider_union_id = ? AND provider = ?", unionID, provider).First(&oa).Error
	return oa, err
}

func (r *Repo) FindOAuthAccountByProviderOpenID(provider common.ProviderType, openID string) (OAuthAccount, error) {
	var oa OAuthAccount
	err := r.db.Where("provider = ? AND provider_open_id = ?", provider, openID).First(&oa).Error
	return oa, err
}

func (r *Repo) CreateOAuthAccount(oa *OAuthAccount) error {
	return r.db.Create(oa).Error
}

