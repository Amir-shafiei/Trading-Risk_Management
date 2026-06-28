package repository

import (
	"Trading-Risk_Management/models"
	"time"

	"gorm.io/gorm"
)

type TradeRepository interface {
	Create(trade *models.Trade) error
	GetByID(id uint) (*models.Trade, error)
	GetByUserID(userID uint, limit, offset int) ([]models.Trade, error)
	Update(trade *models.Trade) error
	Delete(id uint) error
	GetPnLHistory(userID uint) ([]models.Trade, error)
	GetStats(userID uint) (*models.TradeStats, error)
	GetTodayPnL(userID uint) (float64, error)
	GetTodayClosedTrades(userID uint) ([]models.Trade, error)
	CountOpenByPortfolio(portfolioID uint) (int, error)
	GetClosedTrades(userID uint) ([]models.Trade, error)
}
type TradeRepo struct {
	db *gorm.DB
}

func NewTradeRepo(db *gorm.DB) TradeRepository {
	return &TradeRepo{db: db}
}
func (repo *TradeRepo) Create(trade *models.Trade) error {
	return repo.db.Create(trade).Error
}
func (repo *TradeRepo) GetByID(id uint) (*models.Trade, error) {
	var trade models.Trade
	result := repo.db.Where("id = ?", id).First(&trade)
	return &trade, result.Error
}
func (repo *TradeRepo) GetByUserID(userID uint, limit, offset int) ([]models.Trade, error) {
	var trades []models.Trade

	err := repo.db.
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).
		Offset(offset).
		Find(&trades).Error

	return trades, err
}
func (repo *TradeRepo) Update(trade *models.Trade) error {
	return repo.db.Save(trade).Error
}

func (repo *TradeRepo) Delete(id uint) error {
	var trade models.Trade
	result := repo.db.Where("id = ?", id).Delete(&trade)
	return result.Error
}
func (repo *TradeRepo) GetPnLHistory(userID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := repo.db.
		Where("user_id = ? AND status != ?", userID, models.TradeOpen).
		Order("closed_at asc").
		Find(&trades).Error
	return trades, err
}
func (repo *TradeRepo) GetStats(userID uint) (*models.TradeStats, error) {
	var allTrades []models.Trade
	err := repo.db.Where("user_id = ?", userID).Find(&allTrades).Error
	if err != nil {
		return nil, err
	}

	stats := &models.TradeStats{}
	var wins int
	var rrSum float64
	var rrCount int

	for _, t := range allTrades {
		if t.Status == models.TradeOpen {
			stats.OpenTrades++
		} else {
			stats.ClosedTrades++
			if t.PnL != nil {
				stats.TotalPnL += *t.PnL
				if *t.PnL > stats.BestTrade {
					stats.BestTrade = *t.PnL
				}
				if *t.PnL < stats.WorstTrade {
					stats.WorstTrade = *t.PnL
				}
				if *t.PnL > 0 {
					wins++
				}
			}
			if t.TakeProfit != nil && t.EntryPrice > 0 && t.StopLoss > 0 {
				var rr float64
				if t.Side == models.Long {
					rr = (*t.TakeProfit - t.EntryPrice) / (t.EntryPrice - t.StopLoss)
				} else {
					rr = (t.EntryPrice - *t.TakeProfit) / (t.StopLoss - t.EntryPrice)
				}
				if rr > 0 {
					rrSum += rr
					rrCount++
				}
			}
		}
	}

	if stats.ClosedTrades > 0 {
		stats.WinRate = float64(wins) / float64(stats.ClosedTrades) * 100
	}
	if rrCount > 0 {
		stats.AvgRiskReward = rrSum / float64(rrCount)
	}

	return stats, nil
}

func (repo *TradeRepo) GetTodayPnL(userID uint) (float64, error) {
	var total float64
	today := time.Now().Format("2006-01-02")
	err := repo.db.
		Table("trades").
		Select("COALESCE(SUM(pnl), 0)").
		Where("user_id = ? AND status != ? AND DATE(closed_at) = ?", userID, models.TradeOpen, today).
		Scan(&total).Error
	return total, err
}

func (repo *TradeRepo) GetTodayClosedTrades(userID uint) ([]models.Trade, error) {
	var trades []models.Trade
	today := time.Now().Format("2006-01-02")
	err := repo.db.
		Where("user_id = ? AND status != ? AND DATE(closed_at) = ?", userID, models.TradeOpen, today).
		Find(&trades).Error
	return trades, err
}

func (repo *TradeRepo) CountOpenByPortfolio(portfolioID uint) (int, error) {
	var count int64
	err := repo.db.
		Model(&models.Trade{}).
		Where("portfolio_id = ? AND status = ?", portfolioID, models.TradeOpen).
		Count(&count).Error
	return int(count), err
}

func (repo *TradeRepo) GetClosedTrades(userID uint) ([]models.Trade, error) {
	var trades []models.Trade
	err := repo.db.
		Where("user_id = ? AND status != ?", userID, models.TradeOpen).
		Order("closed_at asc").
		Find(&trades).Error
	return trades, err
}
