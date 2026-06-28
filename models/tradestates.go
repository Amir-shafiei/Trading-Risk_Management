package models

type TradeStats struct {
	TotalPnL      float64
	WinRate       float64
	OpenTrades    int
	ClosedTrades  int
	BestTrade     float64
	WorstTrade    float64
	AvgRiskReward float64
}
