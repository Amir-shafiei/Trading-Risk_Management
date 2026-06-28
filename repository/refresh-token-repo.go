package repository

import (
	"Trading-Risk_Management/models"

	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	GetByToken(token string) (*models.RefreshToken, error)
	Revoke(token string) error
	RevokeAllForUser(userID uint) error
}

type RefreshTokenRepo struct {
	db *gorm.DB
}

func NewRefreshTokenRepo(db *gorm.DB) RefreshTokenRepository {
	return &RefreshTokenRepo{db: db}
}

func (repo *RefreshTokenRepo) Create(token *models.RefreshToken) error {
	return repo.db.Create(token).Error
}

func (repo *RefreshTokenRepo) GetByToken(token string) (*models.RefreshToken, error) {
	var rt models.RefreshToken
	result := repo.db.Where("token = ? AND revoked = ?", token, false).First(&rt)
	return &rt, result.Error
}

func (repo *RefreshTokenRepo) Revoke(token string) error {
	return repo.db.Model(&models.RefreshToken{}).Where("token = ?", token).Update("revoked", true).Error
}

func (repo *RefreshTokenRepo) RevokeAllForUser(userID uint) error {
	return repo.db.Model(&models.RefreshToken{}).Where("user_id = ?", userID).Update("revoked", true).Error
}
