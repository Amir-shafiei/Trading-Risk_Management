package repository

import (
	"Trading-Risk_Management/models"
	"errors"

	"gorm.io/gorm"
)

type PortfolioRepository interface {
	Create(portfolio *models.Portfolio) error
	GetByID(id uint) (*models.Portfolio, error)
	GetByUserID(userID uint) ([]models.Portfolio, error)
	GetDefault(userID uint) (*models.Portfolio, error)
	Update(portfolio *models.Portfolio) error
	Delete(id uint) error
}
type PtRepo struct {
	db *gorm.DB
}

func NewPtRepo(db *gorm.DB) PortfolioRepository {
	return &PtRepo{db: db}
}
func (repo *PtRepo) Create(portfolio *models.Portfolio) error {
	return repo.db.Create(portfolio).Error
}
func (repo *PtRepo) GetByID(id uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	result := repo.db.Where("id = ?", id).First(&portfolio)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("portfolio not found")
	}
	return &portfolio, result.Error
}
func (repo *PtRepo) GetByUserID(userID uint) ([]models.Portfolio, error) {
	var portfolios []models.Portfolio
	result := repo.db.Where("user_id = ?", userID).Find(&portfolios)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")

	}
	return portfolios, nil
}
func (repo *PtRepo) GetDefault(userID uint) (*models.Portfolio, error) {
	var portfolio models.Portfolio
	result := repo.db.Where("user_id = ? AND is_default = ?", userID, true).First(&portfolio)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	return &portfolio, nil
}
func (repo *PtRepo) Update(portfolio *models.Portfolio) error {
	return repo.db.Save(portfolio).Error

}
func (repo *PtRepo) Delete(id uint) error {
	var portfolio models.Portfolio
	result := repo.db.Where("id = ?", id).Delete(&portfolio)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("user not found")
	}
	return result.Error
}
