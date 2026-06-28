package repository

import (
	"Trading-Risk_Management/models"

	"gorm.io/gorm"
)

type ChecklistRepository interface {
	Create(checklist *models.PreTradeChecklist) error
	GetByTrade(userID, tradeID uint) (*models.PreTradeChecklist, error)
	Update(checklist *models.PreTradeChecklist) error
	GetDefaults(userID uint) (string, error)
	SetDefaults(userID uint, items string) error
}

type ChecklistRepo struct {
	db *gorm.DB
}

func NewChecklistRepo(db *gorm.DB) ChecklistRepository {
	return &ChecklistRepo{db: db}
}

func (repo *ChecklistRepo) Create(checklist *models.PreTradeChecklist) error {
	return repo.db.Create(checklist).Error
}

func (repo *ChecklistRepo) GetByTrade(userID, tradeID uint) (*models.PreTradeChecklist, error) {
	var checklist models.PreTradeChecklist
	result := repo.db.Where("user_id = ? AND trade_id = ?", userID, tradeID).First(&checklist)
	return &checklist, result.Error
}

func (repo *ChecklistRepo) Update(checklist *models.PreTradeChecklist) error {
	return repo.db.Save(checklist).Error
}

func (repo *ChecklistRepo) GetDefaults(userID uint) (string, error) {
	var defaults models.ChecklistDefaults
	result := repo.db.Where("user_id = ?", userID).First(&defaults)
	if result.Error != nil {
		return "[]", nil
	}
	return defaults.Items, nil
}

func (repo *ChecklistRepo) SetDefaults(userID uint, items string) error {
	var defaults models.ChecklistDefaults
	result := repo.db.Where("user_id = ?", userID).First(&defaults)
	if result.Error != nil {
		defaults = models.ChecklistDefaults{UserID: userID, Items: items}
		return repo.db.Create(&defaults).Error
	}
	defaults.Items = items
	return repo.db.Save(&defaults).Error
}
