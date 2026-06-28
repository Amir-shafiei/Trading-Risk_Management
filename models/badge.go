package models

import "time"

type Badge struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id" gorm:"index"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	EarnedAt    time.Time `json:"earned_at"`
}

type BadgeDefinition struct {
	Name        string
	Description string
	Icon        string
	Check       func(stats *TradeStats, totalTrades int) bool
}

var AllBadges = []BadgeDefinition{
	{
		Name: "\u0627\u0648\u0644\u06cc\u0646 \u062e\u0648\u0646",
		Description: "\u0627\u0648\u0644\u06cc\u0646 \u062a\u0631\u06cc\u062f \u062e\u0648\u062f \u0631\u0627 \u0628\u0628\u0646\u062f\u06cc\u062f",
		Icon: "blood",
		Check: func(s *TradeStats, total int) bool { return s.ClosedTrades >= 1 },
	},
	{
		Name: "\u0634\u0648\u0627\u0639\u062a\u064b \u0628\u0631\u0646\u062f\u0647",
		Description: "\u06f5 \u062a\u0631\u06cc\u062f \u0628\u0631\u0646\u062f\u0647 \u067e\u0634\u062a \u0633\u0631 \u0647\u0645",
		Icon: "fire",
		Check: func(s *TradeStats, total int) bool { return s.WinRate >= 100 && s.ClosedTrades >= 5 },
	},
	{
		Name: "\u062a\u0648\u0642\u0641\u200c\u0646\u0627\u067e\u0630\u06cc\u0631",
		Description: "\u06f1\u06f0 \u062a\u0631\u06cc\u062f \u0628\u0631\u0646\u062f\u0647 \u067e\u0634\u062a \u0633\u0631 \u0647\u0645",
		Icon: "trophy",
		Check: func(s *TradeStats, total int) bool { return s.WinRate >= 100 && s.ClosedTrades >= 10 },
	},
	{
		Name: "\u062b\u0628\u0627\u062a\u200c\u0642\u062f\u0645",
		Description: "\u0646\u0631\u062e \u0628\u0631\u062f \u0628\u0627\u0644\u0627\u06cc \u06f6\u06f0\u066a \u0628\u0627 \u06f2\u06f0+ \u062a\u0631\u06cc\u062f",
		Icon: "chart",
		Check: func(s *TradeStats, total int) bool { return s.WinRate > 60 && s.ClosedTrades >= 20 },
	},
	{
		Name: "\u062a\u06a9\u200c\u062a\u06cc\u0631\u0627\u0646\u062f\u0627\u0632",
		Description: "\u0646\u0631\u062e \u0628\u0631\u062f \u0628\u0627\u0644\u0627\u06cc \u06f7\u06f5\u066a \u0628\u0627 \u06f1\u06f0+ \u062a\u0631\u06cc\u062f",
		Icon: "target",
		Check: func(s *TradeStats, total int) bool { return s.WinRate > 75 && s.ClosedTrades >= 10 },
	},
	{
		Name: "\u0628\u0631\u0646\u062f\u0647 \u0628\u0632\u0631\u06af",
		Description: "\u0633\u0648\u062f \u06cc\u06a9 \u062a\u0631\u06cc\u062f \u0628\u06cc\u0634 \u0627\u0632 \u06f5\u066a \u067e\u0648\u0631\u062a\u0641\u0648\u0644\u06cc\u0648",
		Icon: "money",
		Check: func(s *TradeStats, total int) bool { return s.BestTrade > 0 },
	},
	{
		Name: "\u0645\u062f\u06cc\u0631 \u0631\u06cc\u0633\u06a9",
		Description: "\u06f1\u06f0 \u062a\u0631\u06cc\u062f \u0628\u0627 \u0631\u06cc\u0633\u06a9 \u06a9\u0645\u062a\u0631 \u0627\u0632 \u06f1\u066a",
		Icon: "shield",
		Check: func(s *TradeStats, total int) bool { return s.ClosedTrades >= 10 },
	},
	{
		Name: "\u06a9\u0647\u0646\u0647\u200c\u0633\u0631\u0628\u0627\u0632",
		Description: "\u062a\u06a9\u0645\u06cc\u0644 \u06f5\u06f0 \u062a\u0631\u06cc\u062f",
		Icon: "medal",
		Check: func(s *TradeStats, total int) bool { return s.ClosedTrades >= 50 },
	},
	{
		Name: "\u062f\u0633\u062a\u200c\u0622\u0647\u0646\u06cc\u0646",
		Description: "\u0646\u06af\u0647\u200c\u062f\u0627\u0634\u062a\u0646 \u062a\u0631\u06cc\u062f \u0628\u06cc\u0634 \u0627\u0632 \u06f2\u06f4 \u0633\u0627\u0639\u062a \u0628\u0627 \u0633\u0648\u062f",
		Icon: "hand",
		Check: func(s *TradeStats, total int) bool { return s.ClosedTrades >= 1 && s.TotalPnL > 0 },
	},
	{
		Name: "\u0647\u0641\u062a\u0647 \u0628\u06cc\u200c\u0646\u0642\u0636",
		Description: "\u0628\u0631\u062f\u0646 \u0647\u0645\u0647 \u062a\u0631\u06cc\u062f\u0647\u0627 \u062f\u0631 \u06cc\u06a9 \u0647\u0641\u062a\u0647 (\u062d\u062f\u0627\u0642\u0644 \u06f3)",
		Icon: "star",
		Check: func(s *TradeStats, total int) bool { return s.WinRate >= 100 && s.ClosedTrades >= 3 },
	},
}
