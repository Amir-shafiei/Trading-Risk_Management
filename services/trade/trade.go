package trade

import (
	"Trading-Risk_Management/models"
	"Trading-Risk_Management/repository"
	"errors"
	"math"
	"time"
)

type BadgeChecker interface {
	CheckAndAward(userID uint) ([]models.Badge, error)
}

type TradeService interface {
	CreateTrade(trade *models.Trade) error
	CloseTrade(userID, tradeID uint, exitPrice float64) error
	PartialClose(userID, tradeID uint, percent float64, exitPrice float64) error
	MoveToBreakeven(userID, tradeID uint) error
	UpdateTrade(userID uint, tradeID uint, updates *models.Trade) error
	DeleteTrade(userID, tradeID uint) error
	GetByID(userID, tradeID uint) (*models.Trade, error)
	GetByUserID(userID uint, limit, offset int) ([]models.Trade, error)
}

type TradeServiceImpl struct {
	tradeRepo     repository.TradeRepository
	portfolioRepo repository.PortfolioRepository
	badgeChecker  BadgeChecker
}

func NewTradeService(tradeRepo repository.TradeRepository, portfolioRepo repository.PortfolioRepository, badgeChecker BadgeChecker) TradeService {
	return &TradeServiceImpl{
		tradeRepo:     tradeRepo,
		portfolioRepo: portfolioRepo,
		badgeChecker:  badgeChecker,
	}
}

func (s *TradeServiceImpl) CreateTrade(trade *models.Trade) error {
	portfolio, err := s.portfolioRepo.GetByID(trade.PortfolioID)
	if err != nil {
		return err
	}
	if portfolio.UserID != trade.UserID {
		return errors.New("unauthorized")
	}

	if trade.EntryPrice <= 0 {
		return errors.New("entry price must be greater than zero")
	}
	if trade.StopLoss <= 0 {
		return errors.New("stop loss must be greater than zero")
	}
	if trade.RiskPercent <= 0 || trade.RiskPercent > 100 {
		return errors.New("risk percent must be between 0 and 100")
	}

	if trade.Side == models.Long && trade.StopLoss >= trade.EntryPrice {
		return errors.New("stop loss must be below entry price for long positions")
	}
	if trade.Side == models.Short && trade.StopLoss <= trade.EntryPrice {
		return errors.New("stop loss must be above entry price for short positions")
	}

	if portfolio.MaxDailyLoss > 0 {
		todayPnL, err := s.tradeRepo.GetTodayPnL(trade.UserID)
		if err == nil && todayPnL <= -portfolio.MaxDailyLoss {
			return errors.New("max daily loss limit reached, trading blocked for today")
		}
	}

	if portfolio.MaxOpenTrades > 0 {
		openCount, err := s.tradeRepo.CountOpenByPortfolio(trade.PortfolioID)
		if err == nil && openCount >= portfolio.MaxOpenTrades {
			return errors.New("max open trades limit reached for this portfolio")
		}
	}

	riskAmount := portfolio.Balance * (trade.RiskPercent / 100)
	priceDiff := math.Abs(trade.EntryPrice - trade.StopLoss)
	trade.Quantity = riskAmount / priceDiff
	trade.InitialQty = trade.Quantity
	trade.Status = models.TradeOpen

	return s.tradeRepo.Create(trade)
}

func (s *TradeServiceImpl) CloseTrade(userID, tradeID uint, exitPrice float64) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return err
	}

	if trade.UserID != userID {
		return errors.New("unauthorized")
	}

	if trade.Status != models.TradeOpen {
		return errors.New("trade is already closed")
	}

	var pnl float64
	if trade.Side == models.Long {
		pnl = (exitPrice - trade.EntryPrice) * trade.Quantity
	} else {
		pnl = (trade.EntryPrice - exitPrice) * trade.Quantity
	}

	now := time.Now()
	trade.ExitPrice = &exitPrice
	trade.PnL = &pnl
	trade.ClosedAt = &now

	if trade.TakeProfit != nil && exitPrice >= *trade.TakeProfit {
		trade.Status = models.TradeTPHit
	} else if exitPrice == trade.StopLoss {
		trade.Status = models.TradeSLHit
	} else {
		trade.Status = models.TradeClosed
	}

	portfolio, err := s.portfolioRepo.GetByID(trade.PortfolioID)
	if err != nil {
		return err
	}
	portfolio.Balance += pnl
	portfolio.PnL += pnl

	if err := s.portfolioRepo.Update(portfolio); err != nil {
		return err
	}

	if err := s.tradeRepo.Update(trade); err != nil {
		return err
	}

	if s.badgeChecker != nil {
		s.badgeChecker.CheckAndAward(userID)
	}

	return nil
}

func (s *TradeServiceImpl) PartialClose(userID, tradeID uint, percent float64, exitPrice float64) error {
	if percent <= 0 || percent > 100 {
		return errors.New("percent must be between 0 and 100")
	}

	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return err
	}

	if trade.UserID != userID {
		return errors.New("unauthorized")
	}

	if trade.Status != models.TradeOpen {
		return errors.New("trade is already closed")
	}

	closingQty := trade.Quantity * (percent / 100)

	var pnl float64
	if trade.Side == models.Long {
		pnl = (exitPrice - trade.EntryPrice) * closingQty
	} else {
		pnl = (trade.EntryPrice - exitPrice) * closingQty
	}

	trade.Quantity -= closingQty
	trade.PnL = &pnl

	portfolio, err := s.portfolioRepo.GetByID(trade.PortfolioID)
	if err != nil {
		return err
	}
	portfolio.Balance += pnl
	portfolio.PnL += pnl

	if err := s.portfolioRepo.Update(portfolio); err != nil {
		return err
	}

	if trade.Quantity <= 0 {
		now := time.Now()
		trade.ExitPrice = &exitPrice
		trade.ClosedAt = &now
		trade.Status = models.TradeClosed
	}

	return s.tradeRepo.Update(trade)
}

func (s *TradeServiceImpl) MoveToBreakeven(userID, tradeID uint) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return err
	}

	if trade.UserID != userID {
		return errors.New("unauthorized")
	}

	if trade.Status != models.TradeOpen {
		return errors.New("trade is already closed")
	}

	trade.StopLoss = trade.EntryPrice
	return s.tradeRepo.Update(trade)
}

func (s *TradeServiceImpl) GetByID(userID, tradeID uint) (*models.Trade, error) {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return nil, err
	}

	if trade.UserID != userID {
		return nil, errors.New("unauthorized")
	}

	return trade, nil
}

func (s *TradeServiceImpl) UpdateTrade(userID uint, tradeID uint, updates *models.Trade) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return err
	}

	if trade.UserID != userID {
		return errors.New("unauthorized")
	}

	if trade.Status != models.TradeOpen {
		if updates.Note != "" {
			trade.Note = updates.Note
			return s.tradeRepo.Update(trade)
		}
		return errors.New("only note can be edited on closed trades")
	}

	portfolio, err := s.portfolioRepo.GetByID(trade.PortfolioID)
	if err != nil {
		return err
	}

	if updates.Symbol != "" {
		trade.Symbol = updates.Symbol
	}
	if updates.Side != "" {
		trade.Side = updates.Side
	}
	if updates.EntryPrice > 0 {
		trade.EntryPrice = updates.EntryPrice
	}
	if updates.StopLoss > 0 {
		trade.StopLoss = updates.StopLoss
	}
	if updates.TakeProfit != nil {
		trade.TakeProfit = updates.TakeProfit
	}
	if updates.RiskPercent > 0 && updates.RiskPercent <= 100 {
		trade.RiskPercent = updates.RiskPercent
	}
	if updates.Leverage > 0 {
		trade.Leverage = updates.Leverage
	}
	if updates.Note != "" {
		trade.Note = updates.Note
	}

	if trade.Side == models.Long && trade.StopLoss >= trade.EntryPrice {
		return errors.New("stop loss must be below entry price for long positions")
	}
	if trade.Side == models.Short && trade.StopLoss <= trade.EntryPrice {
		return errors.New("stop loss must be above entry price for short positions")
	}

	riskAmount := portfolio.Balance * (trade.RiskPercent / 100)
	priceDiff := math.Abs(trade.EntryPrice - trade.StopLoss)
	trade.Quantity = riskAmount / priceDiff

	return s.tradeRepo.Update(trade)
}

func (s *TradeServiceImpl) DeleteTrade(userID, tradeID uint) error {
	trade, err := s.tradeRepo.GetByID(tradeID)
	if err != nil {
		return err
	}

	if trade.UserID != userID {
		return errors.New("unauthorized")
	}

	if trade.Status != models.TradeOpen && trade.PnL != nil {
		portfolio, err := s.portfolioRepo.GetByID(trade.PortfolioID)
		if err != nil {
			return err
		}
		portfolio.Balance -= *trade.PnL
		portfolio.PnL -= *trade.PnL
		if err := s.portfolioRepo.Update(portfolio); err != nil {
			return err
		}
	}

	return s.tradeRepo.Delete(tradeID)
}

func (s *TradeServiceImpl) GetByUserID(userID uint, limit, offset int) ([]models.Trade, error) {
	if limit <= 0 {
		limit = 10
	}
	return s.tradeRepo.GetByUserID(userID, limit, offset)
}
