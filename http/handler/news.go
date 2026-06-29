package handler

import (
	"Trading-Risk_Management/services/news"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewsHandler struct {
	newsService news.NewsService
}

func NewNewsHandler(s news.NewsService) *NewsHandler {
	return &NewsHandler{newsService: s}
}

func (h *NewsHandler) GetNews(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	events := h.newsService.GetLatest()
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (h *NewsHandler) RefreshNews(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	h.newsService.Refresh()
	c.JSON(http.StatusOK, gin.H{"message": "calendar refreshed"})
}
