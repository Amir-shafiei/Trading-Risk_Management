package models

import "time"

type DailyStreak struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"index;uniqueIndex:idx_user_date"`
	Date      time.Time `json:"date" gorm:"type:date;uniqueIndex:idx_user_date"`
	Wins      int       `json:"wins"`
	Losses    int       `json:"losses"`
	ConsecWin int       `json:"consec_win"`
}
