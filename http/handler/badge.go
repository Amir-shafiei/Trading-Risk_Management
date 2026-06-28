package handler

import (
	"Trading-Risk_Management/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BadgeService interface {
	CheckAndAward(userID uint) ([]models.Badge, error)
	GetUserBadges(userID uint) ([]models.Badge, error)
}

type BadgeHandler struct {
	badgeService BadgeService
}

func NewBadgeHandler(s BadgeService) *BadgeHandler {
	return &BadgeHandler{badgeService: s}
}

func (h *BadgeHandler) CheckBadges(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	newBadges, err := h.badgeService.CheckAndAward(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"new_badges": newBadges,
		"count":      len(newBadges),
	})
}

func (h *BadgeHandler) GetBadges(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	badges, err := h.badgeService.GetUserBadges(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"badges": badges})
}
