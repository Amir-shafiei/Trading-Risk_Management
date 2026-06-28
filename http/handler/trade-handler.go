package handler

import (
	"Trading-Risk_Management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

type TradeHandler struct {
	tradeService TradeService
}

func NewTradeHandler(s TradeService) *TradeHandler {
	return &TradeHandler{tradeService: s}
}

func (h *TradeHandler) CreateTrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var trade models.Trade
	if err := c.ShouldBindJSON(&trade); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	trade.UserID = userID.(uint)

	if err := h.tradeService.CreateTrade(&trade); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "trade created", "trade": trade})
}

func (h *TradeHandler) CloseTrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var body struct {
		ExitPrice float64 `json:"exit_price"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.ExitPrice <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exit price must be greater than zero"})
		return
	}

	if err := h.tradeService.CloseTrade(userID.(uint), uint(tradeID), body.ExitPrice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trade closed"})
}

func (h *TradeHandler) PartialClose(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var body struct {
		Percent   float64 `json:"percent"`
		ExitPrice float64 `json:"exit_price"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Percent <= 0 || body.Percent > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "percent must be between 0 and 100"})
		return
	}
	if body.ExitPrice <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "exit price must be greater than zero"})
		return
	}

	if err := h.tradeService.PartialClose(userID.(uint), uint(tradeID), body.Percent, body.ExitPrice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trade partially closed"})
}

func (h *TradeHandler) MoveToBreakeven(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	if err := h.tradeService.MoveToBreakeven(userID.(uint), uint(tradeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stop loss moved to breakeven"})
}

func (h *TradeHandler) GetTrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	trade, err := h.tradeService.GetByID(userID.(uint), uint(tradeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trade": trade})
}

func (h *TradeHandler) GetTrades(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	trades, err := h.tradeService.GetByUserID(userID.(uint), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"trades": trades})
}

func (h *TradeHandler) UpdateTrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var updates models.Trade
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.tradeService.UpdateTrade(userID.(uint), uint(tradeID), &updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trade updated"})
}

func (h *TradeHandler) DeleteTrade(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	if err := h.tradeService.DeleteTrade(userID.(uint), uint(tradeID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "trade deleted"})
}
