package repository

import (
	"Trading-Risk_Management/models"

	"gorm.io/gorm"
)

type BadgeRepository interface {
	Create(badge *models.Badge) error
	GetByUserID(userID uint) ([]models.Badge, error)
}

type BadgeRepo struct {
	db *gorm.DB
}

func NewBadgeRepo(db *gorm.DB) BadgeRepository {
	return &BadgeRepo{db: db}
}

func (repo *BadgeRepo) Create(badge *models.Badge) error {
	return repo.db.Create(badge).Error
}

func (repo *BadgeRepo) GetByUserID(userID uint) ([]models.Badge, error) {
	var badges []models.Badge
	err := repo.db.Where("user_id = ?", userID).Order("earned_at desc").Find(&badges).Error
	return badges, err
}
