package models

import (
	"time"

	"gorm.io/gorm"
)

type Trade struct {
	gorm.Model
	UserID      uint         `json:"user_id"`
	PortfolioID uint         `json:"portfolio_id"`
	Symbol      string       `json:"symbol"`
	Side        PositionSide `json:"side"`
	EntryPrice  float64      `json:"entry_price" gorm:"column:entry_price"`
	ExitPrice   *float64     `json:"exit_price,omitempty" gorm:"column:exit_price"`
	StopLoss    float64      `json:"stop_loss" gorm:"column:stop_loss"`
	TakeProfit  *float64     `json:"take_profit,omitempty" gorm:"column:take_profit"`
	Leverage    float64      `json:"leverage"`
	RiskPercent float64      `json:"risk_percent" gorm:"column:risk_percent"`
	Quantity    float64      `json:"quantity"`
	InitialQty  float64      `json:"initial_qty" gorm:"column:initial_qty"`
	PnL         *float64     `json:"pnl,omitempty" gorm:"column:pnl"`
	Note        string       `json:"note,omitempty"`
	ClosedAt    *time.Time   `json:"closed_at,omitempty" gorm:"column:closed_at"`
	Status      TradeStatus  `json:"status"`
	Portfolio   Portfolio    `json:"-" gorm:"foreignKey:PortfolioID"`
}

type PositionSide string

const (
	Long  PositionSide = "LONG"
	Short PositionSide = "SHORT"
)

type TradeStatus string

const (
	TradeOpen   TradeStatus = "OPEN"
	TradeClosed TradeStatus = "CLOSED"
	TradeTPHit  TradeStatus = "TP_HIT"
	TradeSLHit  TradeStatus = "SL_HIT"
)
