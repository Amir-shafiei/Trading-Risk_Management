package handler

import (
	"Trading-Risk_Management/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PortfolioSettingsHandler struct {
	portfolioRepo repository.PortfolioRepository
}

func NewPortfolioSettingsHandler(pr repository.PortfolioRepository) *PortfolioSettingsHandler {
	return &PortfolioSettingsHandler{portfolioRepo: pr}
}

func (h *PortfolioSettingsHandler) SetDailyLossLimit(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	portfolioID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid portfolio id"})
		return
	}

	portfolio, err := h.portfolioRepo.GetByID(uint(portfolioID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "portfolio not found"})
		return
	}
	if portfolio.UserID != userID.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		MaxDailyLoss float64 `json:"max_daily_loss"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio.MaxDailyLoss = body.MaxDailyLoss
	if err := h.portfolioRepo.Update(portfolio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "daily loss limit updated", "max_daily_loss": portfolio.MaxDailyLoss})
}

func (h *PortfolioSettingsHandler) SetMaxOpenTrades(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	portfolioID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid portfolio id"})
		return
	}

	portfolio, err := h.portfolioRepo.GetByID(uint(portfolioID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "portfolio not found"})
		return
	}
	if portfolio.UserID != userID.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		MaxOpenTrades int `json:"max_open_trades"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	portfolio.MaxOpenTrades = body.MaxOpenTrades
	if err := h.portfolioRepo.Update(portfolio); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "max open trades updated", "max_open_trades": portfolio.MaxOpenTrades})
}

func (h *PortfolioSettingsHandler) GetDailyLossStatus(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	portfolio, err := h.portfolioRepo.GetDefault(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "portfolio not found"})
		return
	}

	_ = portfolio

	type DailyLossStatus struct {
		MaxDailyLoss float64 `json:"max_daily_loss"`
		CurrentLoss  float64 `json:"current_loss"`
		LimitReached bool    `json:"limit_reached"`
		Remaining    float64 `json:"remaining"`
	}

	status := DailyLossStatus{
		MaxDailyLoss: portfolio.MaxDailyLoss,
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}
