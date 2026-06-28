package handler

import (
	"Trading-Risk_Management/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AlertService interface {
	CheckAlerts(userID uint) ([]models.RiskAlert, error)
	GetUnread(userID uint) ([]models.RiskAlert, error)
	MarkRead(userID, alertID uint) error
}

type AlertHandler struct {
	alertService AlertService
}

func NewAlertHandler(s AlertService) *AlertHandler {
	return &AlertHandler{alertService: s}
}

func (h *AlertHandler) CheckAlerts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	alerts, err := h.alertService.CheckAlerts(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

func (h *AlertHandler) GetUnread(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	alerts, err := h.alertService.GetUnread(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"alerts": alerts})
}

func (h *AlertHandler) MarkRead(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	alertID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid alert id"})
		return
	}

	if err := h.alertService.MarkRead(userID.(uint), uint(alertID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "alert marked as read"})
}
