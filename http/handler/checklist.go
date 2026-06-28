package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChecklistService interface {
	Create(userID, tradeID uint, items []string) error
	GetByTrade(userID, tradeID uint) (interface{}, error)
	UpdateCheck(userID, tradeID uint, itemIndex int, checked bool) error
	GetDefaults(userID uint) ([]string, error)
	SetDefaults(userID uint, items []string) error
}

type ChecklistHandler struct {
	checklistService ChecklistService
}

func NewChecklistHandler(s ChecklistService) *ChecklistHandler {
	return &ChecklistHandler{checklistService: s}
}

func (h *ChecklistHandler) CreateChecklist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		TradeID uint     `json:"trade_id"`
		Items   []string `json:"items"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.checklistService.Create(userID.(uint), body.TradeID, body.Items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "checklist created"})
}

func (h *ChecklistHandler) GetChecklist(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("tradeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	checklist, err := h.checklistService.GetByTrade(userID.(uint), uint(tradeID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "checklist not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"checklist": checklist})
}

func (h *ChecklistHandler) UpdateCheck(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tradeID, err := strconv.ParseUint(c.Param("tradeId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trade id"})
		return
	}

	var body struct {
		ItemIndex int  `json:"item_index"`
		Checked   bool `json:"checked"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.checklistService.UpdateCheck(userID.(uint), uint(tradeID), body.ItemIndex, body.Checked); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "checklist updated"})
}

func (h *ChecklistHandler) GetDefaults(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	items, err := h.checklistService.GetDefaults(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"defaults": items})
}

func (h *ChecklistHandler) SetDefaults(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var body struct {
		Items []string `json:"items"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.checklistService.SetDefaults(userID.(uint), body.Items); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "defaults updated"})
}
