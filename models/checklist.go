package models

import "gorm.io/gorm"

type PreTradeChecklist struct {
	gorm.Model
	UserID  uint   `json:"user_id" gorm:"index"`
	TradeID uint   `json:"trade_id" gorm:"index"`
	Items   string `json:"items"`
	AllMet  bool   `json:"all_met" gorm:"default:false"`
}

type ChecklistDefaults struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"uniqueIndex"`
	Items  string `json:"items"`
}
