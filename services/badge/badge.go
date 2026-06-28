package badge

import (
	"Trading-Risk_Management/models"
	"Trading-Risk_Management/repository"
	"time"
)

type BadgeService interface {
	CheckAndAward(userID uint) ([]models.Badge, error)
	GetUserBadges(userID uint) ([]models.Badge, error)
}

type BadgeServiceImpl struct {
	badgeRepo repository.BadgeRepository
	tradeRepo repository.TradeRepository
}

func NewBadgeService(badgeRepo repository.BadgeRepository, tradeRepo repository.TradeRepository) BadgeService {
	return &BadgeServiceImpl{badgeRepo: badgeRepo, tradeRepo: tradeRepo}
}

func (s *BadgeServiceImpl) CheckAndAward(userID uint) ([]models.Badge, error) {
	stats, err := s.tradeRepo.GetStats(userID)
	if err != nil {
		return nil, err
	}

	trades, err := s.tradeRepo.GetClosedTrades(userID)
	if err != nil {
		return nil, err
	}

	existing, _ := s.badgeRepo.GetByUserID(userID)
	existingMap := make(map[string]bool)
	for _, b := range existing {
		existingMap[b.Name] = true
	}

	var newBadges []models.Badge

	for _, def := range models.AllBadges {
		if existingMap[def.Name] {
			continue
		}
		if def.Check(stats, len(trades)) {
			badge := models.Badge{
				UserID:      userID,
				Name:        def.Name,
				Description: def.Description,
				Icon:        def.Icon,
				EarnedAt:    time.Now(),
			}
			s.badgeRepo.Create(&badge)
			newBadges = append(newBadges, badge)
		}
	}

	return newBadges, nil
}

func (s *BadgeServiceImpl) GetUserBadges(userID uint) ([]models.Badge, error) {
	return s.badgeRepo.GetByUserID(userID)
}
