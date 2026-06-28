package models

import "gorm.io/gorm"

type Portfolio struct {
	gorm.Model
	UserID        uint    `json:"user_id"`
	Name          string  `json:"name"`
	Capital       float64 `json:"capital"`
	Balance       float64 `json:"balance"`
	PnL           float64 `json:"pnl"`
	IsDefault     bool    `json:"is_default"`
	MaxDailyLoss  float64 `json:"max_daily_loss" gorm:"default:0"`
	MaxOpenTrades int     `json:"max_open_trades" gorm:"default:0"`
	Trades        []Trade `json:"-" gorm:"foreignKey:PortfolioID"`
}
