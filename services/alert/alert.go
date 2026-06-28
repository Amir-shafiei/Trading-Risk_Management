package alert

import (
	"Trading-Risk_Management/models"
	"Trading-Risk_Management/repository"
	"errors"
	"fmt"
	"time"
)

type AlertService interface {
	CheckAlerts(userID uint) ([]models.RiskAlert, error)
	GetUnread(userID uint) ([]models.RiskAlert, error)
	MarkRead(userID, alertID uint) error
}

type AlertServiceImpl struct {
	tradeRepo     repository.TradeRepository
	portfolioRepo repository.PortfolioRepository
	alertRepo     repository.AlertRepository
}

func NewAlertService(tradeRepo repository.TradeRepository, portfolioRepo repository.PortfolioRepository, alertRepo repository.AlertRepository) AlertService {
	return &AlertServiceImpl{
		tradeRepo:     tradeRepo,
		portfolioRepo: portfolioRepo,
		alertRepo:     alertRepo,
	}
}

func (s *AlertServiceImpl) CheckAlerts(userID uint) ([]models.RiskAlert, error) {
	portfolio, err := s.portfolioRepo.GetDefault(userID)
	if err != nil {
		return nil, err
	}

	var alerts []models.RiskAlert

	if portfolio.MaxDailyLoss > 0 {
		todayPnL, err := s.tradeRepo.GetTodayPnL(userID)
		if err == nil {
			lossRatio := 0.0
			if portfolio.Balance > 0 {
				lossRatio = (-todayPnL / portfolio.Balance) * 100
			}
			thresholdPct := (portfolio.MaxDailyLoss / portfolio.Balance) * 100
			if lossRatio >= thresholdPct*0.8 && lossRatio < thresholdPct {
				alerts = append(alerts, models.RiskAlert{
					UserID:  userID,
					Type:    "daily_loss_warning",
					Message: fmt.Sprintf("\u062a\u0648\u062c\u0647: \u0636\u0631\u0631 \u0631\u0648\u0632\u0627\u0646\u0647 \u0628\u0647 %.1f%%\u060c \u0633\u0631\u062d %.1f%% \u0627\u0633\u062a", lossRatio, thresholdPct),
					Level:   "warning",
				})
			}
			if lossRatio >= thresholdPct {
				alerts = append(alerts, models.RiskAlert{
					UserID:  userID,
					Type:    "daily_loss_breached",
					Message: fmt.Sprintf("\u0635\u0641\u062d\u0631\u0638\u0631\u0633\u0627\u0646\u06cc \u0631\u0648\u0632\u0627\u0646\u0647 \u0634\u06a9\u0633\u062a! \u062a\u0631\u0627\u06a9\u0646\u0634 \u0645\u0633\u062f\u0648\u062f. \u0636\u0631\u0631: %.1f%%\u060c \u0633\u0631\u062d: %.1f%%", lossRatio, thresholdPct),
					Level:   "danger",
				})
			}
		}
	}

	if portfolio.MaxOpenTrades > 0 {
		openCount, err := s.tradeRepo.CountOpenByPortfolio(portfolio.ID)
		if err == nil {
			if openCount >= portfolio.MaxOpenTrades-1 && openCount < portfolio.MaxOpenTrades {
				alerts = append(alerts, models.RiskAlert{
					UserID:  userID,
					Type:    "open_trades_warning",
					Message: fmt.Sprintf("\u062a\u0648\u062c\u0647: %d/%d \u067e\u0648\u0632\u06cc\u0635\u0648\u0646 \u0628\u0627\u0632\u060c \u06cc\u06a9 \u062a\u0631\u06cc\u062f \u062f\u06cc\u06af\u0631 \u0627\u0632 \u0633\u0631\u062d \u0628\u0631\u062e\u0648\u0627\u0647\u062f \u0631\u0633\u06cc\u062f", openCount, portfolio.MaxOpenTrades),
					Level:   "warning",
				})
			}
			if openCount >= portfolio.MaxOpenTrades {
				alerts = append(alerts, models.RiskAlert{
					UserID:  userID,
					Type:    "open_trades_breached",
					Message: fmt.Sprintf("\u062d\u062f \u0627\u06a9\u0633\u0631 \u067e\u0648\u0632\u06cc\u0635\u0648\u0646 \u0628\u0627\u0632 \u0634\u062f! %d/%d \u067e\u0648\u0632\u06cc\u0635\u0648\u0646 \u0628\u0627\u0632 \u0627\u0633\u062a", openCount, portfolio.MaxOpenTrades),
					Level:   "danger",
				})
			}
		}
	}

	todayPnL, _ := s.tradeRepo.GetTodayPnL(userID)
	if todayPnL < 0 && portfolio.Balance > 0 {
		pctLoss := (-todayPnL / portfolio.Balance) * 100
		if pctLoss >= 2 {
			alerts = append(alerts, models.RiskAlert{
				UserID:  userID,
				Type:    "significant_loss",
				Message: fmt.Sprintf("\u0627\u0645\u0631\u0648\u0632 \u06f2\u066a \u067e\u0648\u0631\u062a\u0641\u0648\u0644\u06cc\u0648 \u062e\u0648\u062f \u0631\u0627 \u0627\u0633\u062a. \u0628\u0647\u062a\u0631 \u0627\u0633\u062a \u0627\u0645\u0631\u0648\u0632 \u0631\u0648 \u0642\u0637\u0639 \u06a9\u0646\u06cc\u062f.", pctLoss),
				Level:   "info",
			})
		}
	}

	for i := range alerts {
		alerts[i].CreatedAt = time.Now()
		s.alertRepo.Create(&alerts[i])
	}

	return alerts, nil
}

func (s *AlertServiceImpl) GetUnread(userID uint) ([]models.RiskAlert, error) {
	return s.alertRepo.GetUnread(userID)
}

func (s *AlertServiceImpl) MarkRead(userID, alertID uint) error {
	alert, err := s.alertRepo.GetByID(alertID)
	if err != nil {
		return err
	}
	if alert.UserID != userID {
		return errors.New("unauthorized")
	}
	return s.alertRepo.MarkRead(alertID)
}
