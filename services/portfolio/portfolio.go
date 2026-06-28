package portfolio

import (
	"Trading-Risk_Management/models"
	"Trading-Risk_Management/repository"
	"errors"
)

type PortfolioService interface {
	Create(userID uint, name string, capital float64) error
	GetByUserID(userID uint) ([]models.Portfolio, error)
	GetByID(id uint) (*models.Portfolio, error)
	SetDefault(userID uint, portfolioID uint) error
	Delete(userID uint, portfolioID uint) error
}

type PortfolioServiceImpl struct {
	portfolioRepo repository.PortfolioRepository
}

func NewPortfolioService(portfolioRepo repository.PortfolioRepository) PortfolioService {
	return &PortfolioServiceImpl{portfolioRepo: portfolioRepo}
}

func (s *PortfolioServiceImpl) Create(userID uint, name string, capital float64) error {
	// ۱. چک کن capital معتبره
	if capital <= 0 {
		return errors.New("capital must be greater than zero")
	}

	// ۲. چک کن پورتفولیو پیشفرض داره یا نه
	existing, _ := s.portfolioRepo.GetDefault(userID)
	isDefault := existing == nil // اگه نداشت، این پیشفرض بشه

	portfolio := &models.Portfolio{
		UserID:    userID,
		Name:      name,
		Capital:   capital,
		Balance:   capital, // اول balance برابر capital هست
		IsDefault: isDefault,
	}

	return s.portfolioRepo.Create(portfolio)
}

func (s *PortfolioServiceImpl) GetByUserID(userID uint) ([]models.Portfolio, error) {
	portfolios, err := s.portfolioRepo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}
	if len(portfolios) == 0 {
		return nil, errors.New("no portfolios found")
	}
	return portfolios, nil
}

func (s *PortfolioServiceImpl) GetByID(id uint) (*models.Portfolio, error) {
	return s.portfolioRepo.GetByID(id)
}

func (s *PortfolioServiceImpl) SetDefault(userID uint, portfolioID uint) error {

	current, err := s.portfolioRepo.GetDefault(userID)
	if err == nil && current != nil {
		// ۲. پیشفرض بودنش رو بردار
		current.IsDefault = false
		if err := s.portfolioRepo.Update(current); err != nil {
			return err
		}
	}

	portfolio, err := s.portfolioRepo.GetByID(portfolioID)
	if err != nil {
		return err
	}

	// ۴. چک کن متعلق به همین کاربره
	if portfolio.UserID != userID {
		return errors.New("unauthorized")
	}

	// ۵. پیشفرض کن
	portfolio.IsDefault = true
	return s.portfolioRepo.Update(portfolio)
}

func (s *PortfolioServiceImpl) Delete(userID uint, portfolioID uint) error {
	portfolio, err := s.portfolioRepo.GetByID(portfolioID)
	if err != nil {
		return err
	}

	if portfolio.UserID != userID {
		return errors.New("unauthorized")
	}

	if portfolio.IsDefault {
		return errors.New("cannot delete default portfolio")
	}

	return s.portfolioRepo.Delete(portfolioID)
}
