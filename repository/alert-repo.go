package repository

import (
	"Trading-Risk_Management/models"

	"gorm.io/gorm"
)

type AlertRepository interface {
	Create(alert *models.RiskAlert) error
	GetByID(id uint) (*models.RiskAlert, error)
	GetUnread(userID uint) ([]models.RiskAlert, error)
	MarkRead(id uint) error
}

type AlertRepo struct {
	db *gorm.DB
}

func NewAlertRepo(db *gorm.DB) AlertRepository {
	return &AlertRepo{db: db}
}

func (repo *AlertRepo) Create(alert *models.RiskAlert) error {
	return repo.db.Create(alert).Error
}

func (repo *AlertRepo) GetByID(id uint) (*models.RiskAlert, error) {
	var alert models.RiskAlert
	result := repo.db.Where("id = ?", id).First(&alert)
	return &alert, result.Error
}

func (repo *AlertRepo) GetUnread(userID uint) ([]models.RiskAlert, error) {
	var alerts []models.RiskAlert
	err := repo.db.Where("user_id = ? AND read = ?", userID, false).Order("created_at desc").Find(&alerts).Error
	return alerts, err
}

func (repo *AlertRepo) MarkRead(id uint) error {
	return repo.db.Model(&models.RiskAlert{}).Where("id = ?", id).Update("read", true).Error
}
