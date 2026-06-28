package handler

import (
	"Trading-Risk_Management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PortfolioService interface {
	Create(userID uint, name string, capital float64) error
	GetByUserID(userID uint) ([]models.Portfolio, error)
	GetByID(id uint) (*models.Portfolio, error)
	SetDefault(userID uint, portfolioID uint) error
	Delete(userID uint, portfolioID uint) error
}

type PortfolioHandler struct {
	portfolioService PortfolioService
}

func NewPortfolioHandler(s PortfolioService) *PortfolioHandler {
	return &PortfolioHandler{portfolioService: s}
}

func (h *PortfolioHandler) CreatePortfolio(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Name    string  `json:"name"`
		Capital float64 `json:"capital"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	if err := h.portfolioService.Create(userID.(uint), body.Name, body.Capital); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "portfolio created"})
}

func (h *PortfolioHandler) GetPortfolios(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	portfolios, err := h.portfolioService.GetByUserID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"portfolios": portfolios})
}

func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
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

	portfolio, err := h.portfolioService.GetByID(uint(portfolioID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// چک کن متعلق به همین کاربره
	if portfolio.UserID != userID.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"portfolio": portfolio})
}

func (h *PortfolioHandler) SetDefault(c *gin.Context) {
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

	if err := h.portfolioService.SetDefault(userID.(uint), uint(portfolioID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "default portfolio updated"})
}

func (h *PortfolioHandler) DeletePortfolio(c *gin.Context) {
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

	if err := h.portfolioService.Delete(userID.(uint), uint(portfolioID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "portfolio deleted"})
}
