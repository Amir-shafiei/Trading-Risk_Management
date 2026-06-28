package handler

import (
	"Trading-Risk_Management/services/dashboard"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardService interface {
	GetDashboard(userID uint) (*dashboard.DashboardResponse, error)
	GetPnLHistory(userID uint) ([]dashboard.PnLPoint, error)
	GetDailyPnL(userID uint) ([]dashboard.DailyPnL, error)
	Calculate(portfolioID uint, entryPrice, stopLoss float64, takeProfit *float64, riskPercent float64) (*dashboard.CalculatorResponse, error)
}

type DashboardHandler struct {
	dashboardService DashboardService
}

func NewDashboardHandler(s DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: s}
}

func (h *DashboardHandler) GetDashboard(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.dashboardService.GetDashboard(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboard": result})
}

func (h *DashboardHandler) GetPnLHistory(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.dashboardService.GetPnLHistory(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"pnl_history": result})
}

func (h *DashboardHandler) Calculate(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		PortfolioID uint     `json:"portfolio_id"`
		EntryPrice  float64  `json:"entry_price"`
		StopLoss    float64  `json:"stop_loss"`
		TakeProfit  *float64 `json:"take_profit"`
		RiskPercent float64  `json:"risk_percent"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// اگه portfolio_id نداد، پیشفرض رو بگیر
	if body.PortfolioID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "portfolio_id is required"})
		return
	}

	_ = userID // برای چک امنیتی بعداً میشه استفاده کرد

	result, err := h.dashboardService.Calculate(
		body.PortfolioID,
		body.EntryPrice,
		body.StopLoss,
		body.TakeProfit,
		body.RiskPercent,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

func (h *DashboardHandler) GetDailyPnL(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	result, err := h.dashboardService.GetDailyPnL(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"daily_pnl": result})
}
