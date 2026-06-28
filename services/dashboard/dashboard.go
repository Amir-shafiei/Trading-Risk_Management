package dashboard

import (
	"Trading-Risk_Management/repository"
	"errors"
	"math"
	"time"
)

type DashboardResponse struct {
	Balance        float64 `json:"balance"`
	TotalPnL       float64 `json:"total_pnl"`
	OpenTrades     int     `json:"open_trades"`
	ClosedTrades   int     `json:"closed_trades"`
	WinRate        float64 `json:"win_rate"`
	BestTrade      float64 `json:"best_trade"`
	WorstTrade     float64 `json:"worst_trade"`
	AvgRiskReward  float64 `json:"avg_risk_reward"`
}

type PnLPoint struct {
	Date   time.Time `json:"date"`
	PnL    float64   `json:"pnl"`
	Symbol string    `json:"symbol"`
	Side   string    `json:"side"`
}

type DailyPnL struct {
	Date string  `json:"date"`
	PnL  float64 `json:"pnl"`
}

type CalculatorResponse struct {
	Quantity      float64 `json:"quantity"`
	RiskAmount    float64 `json:"risk_amount"`
	PositionValue float64 `json:"position_value"`
	RiskReward    float64 `json:"risk_reward"`
}

type DashboardService interface {
	GetDashboard(userID uint) (*DashboardResponse, error)
	GetPnLHistory(userID uint) ([]PnLPoint, error)
	GetDailyPnL(userID uint) ([]DailyPnL, error)
	Calculate(portfolioID uint, entryPrice, stopLoss float64, takeProfit *float64, riskPercent float64) (*CalculatorResponse, error)
}

type DashboardServiceImpl struct {
	tradeRepo     repository.TradeRepository
	portfolioRepo repository.PortfolioRepository
}

func NewDashboardService(tradeRepo repository.TradeRepository, portfolioRepo repository.PortfolioRepository) DashboardService {
	return &DashboardServiceImpl{
		tradeRepo:     tradeRepo,
		portfolioRepo: portfolioRepo,
	}
}

func (s *DashboardServiceImpl) GetDashboard(userID uint) (*DashboardResponse, error) {
	// ۱. گرفتن پورتفولیو پیشفرض
	portfolio, err := s.portfolioRepo.GetDefault(userID)
	if err != nil {
		return nil, err
	}

	// ۲. گرفتن stats
	stats, err := s.tradeRepo.GetStats(userID)
	if err != nil {
		return nil, err
	}

	return &DashboardResponse{
		Balance:       portfolio.Balance,
		TotalPnL:      stats.TotalPnL,
		OpenTrades:    stats.OpenTrades,
		ClosedTrades:  stats.ClosedTrades,
		WinRate:       stats.WinRate,
		BestTrade:     stats.BestTrade,
		WorstTrade:    stats.WorstTrade,
		AvgRiskReward: stats.AvgRiskReward,
	}, nil
}

func (s *DashboardServiceImpl) GetPnLHistory(userID uint) ([]PnLPoint, error) {
	// گرفتن trade های بسته شده
	trades, err := s.tradeRepo.GetPnLHistory(userID)
	if err != nil {
		return nil, err
	}

	if len(trades) == 0 {
		return nil, errors.New("no closed trades found")
	}

	// تبدیل به PnLPoint
	points := make([]PnLPoint, 0, len(trades))
	var cumulative float64

	for _, t := range trades {
		if t.PnL == nil || t.ClosedAt == nil {
			continue
		}
		cumulative += *t.PnL
		points = append(points, PnLPoint{
			Date:   *t.ClosedAt,
			PnL:    cumulative, // cumulative pnl برای نمودار
			Symbol: t.Symbol,
			Side:   string(t.Side),
		})
	}

	return points, nil
}

func (s *DashboardServiceImpl) Calculate(portfolioID uint, entryPrice, stopLoss float64, takeProfit *float64, riskPercent float64) (*CalculatorResponse, error) {
	// ۱. گرفتن پورتفولیو
	portfolio, err := s.portfolioRepo.GetByID(portfolioID)
	if err != nil {
		return nil, err
	}

	// ۲. validate
	if entryPrice <= 0 {
		return nil, errors.New("entry price must be greater than zero")
	}
	if stopLoss <= 0 {
		return nil, errors.New("stop loss must be greater than zero")
	}
	if riskPercent <= 0 || riskPercent > 100 {
		return nil, errors.New("risk percent must be between 0 and 100")
	}

	// ۳. محاسبه
	riskAmount := portfolio.Balance * (riskPercent / 100)
	priceDiff := math.Abs(entryPrice - stopLoss)
	if priceDiff == 0 {
		return nil, errors.New("entry price and stop loss cannot be equal")
	}

	quantity := riskAmount / priceDiff
	positionValue := quantity * entryPrice

	// ۴. محاسبه risk/reward اگه takeProfit داشت
	var riskReward float64
	if takeProfit != nil {
		reward := math.Abs(*takeProfit - entryPrice)
		riskReward = reward / priceDiff
	}

	return &CalculatorResponse{
		Quantity:      math.Round(quantity*100000) / 100000,
		RiskAmount:    math.Round(riskAmount*100) / 100,
		PositionValue: math.Round(positionValue*100) / 100,
		RiskReward:    math.Round(riskReward*100) / 100,
	}, nil
}

func (s *DashboardServiceImpl) GetDailyPnL(userID uint) ([]DailyPnL, error) {
	trades, err := s.tradeRepo.GetPnLHistory(userID)
	if err != nil {
		return nil, err
	}

	dailyMap := make(map[string]float64)
	for _, t := range trades {
		if t.PnL == nil || t.ClosedAt == nil {
			continue
		}
		dateKey := t.ClosedAt.Format("2006-01-02")
		dailyMap[dateKey] += *t.PnL
	}

	result := make([]DailyPnL, 0, len(dailyMap))
	for date, pnl := range dailyMap {
		result = append(result, DailyPnL{Date: date, PnL: math.Round(pnl*100) / 100})
	}

	return result, nil
}
