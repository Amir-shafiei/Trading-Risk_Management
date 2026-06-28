package router

import (
	"Trading-Risk_Management/config"
	"Trading-Risk_Management/http/handler"
	"Trading-Risk_Management/http/middleware"
	"Trading-Risk_Management/repository"

	"github.com/gin-gonic/gin"
)

func Setup(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	tradeHandler *handler.TradeHandler,
	portfolioHandler *handler.PortfolioHandler,
	dashboardHandler *handler.DashboardHandler,
	alertHandler *handler.AlertHandler,
	checklistHandler *handler.ChecklistHandler,
	badgeHandler *handler.BadgeHandler,
	portfolioSettingsHandler *handler.PortfolioSettingsHandler,
	newsHandler *handler.NewsHandler,
	portfolioRepo repository.PortfolioRepository,
) *gin.Engine {
	r := gin.Default()

	r.Static("/css", "./frontend/css")
	r.Static("/js", "./frontend/js")

	r.StaticFile("/", "./frontend/pages/landing.html")
	r.StaticFile("/login", "./frontend/pages/login.html")
	r.StaticFile("/register", "./frontend/pages/register.html")
	r.StaticFile("/dashboard", "./frontend/pages/dashboard.html")
	r.StaticFile("/trades", "./frontend/pages/trade.html")
	r.StaticFile("/portfolio", "./frontend/pages/portfolio.html")
	r.StaticFile("/journal", "./frontend/pages/journal.html")
	r.StaticFile("/badges", "./frontend/pages/badges.html")
	r.StaticFile("/news", "./frontend/pages/news.html")
	r.StaticFile("/settings", "./frontend/pages/settings.html")

	public := r.Group("/api")
	{
		public.POST("/register", authHandler.RegHandler)
		public.POST("/login", authHandler.LoginHandler)
		public.POST("/refresh", authHandler.RefreshTokenHandler)
	}

	protected := r.Group("/api")
	protected.Use(middleware.NewAuthMiddleware(cfg))
	{
		protected.POST("/logout", authHandler.LogoutHandler)
		protected.PUT("/user/password", authHandler.ChangePassword)

		protected.POST("/portfolio", portfolioHandler.CreatePortfolio)
		protected.GET("/portfolio", portfolioHandler.GetPortfolios)
		protected.GET("/portfolio/:id", portfolioHandler.GetPortfolio)
		protected.PUT("/portfolio/:id/default", portfolioHandler.SetDefault)
		protected.DELETE("/portfolio/:id", portfolioHandler.DeletePortfolio)
		protected.PUT("/portfolio/:id/daily-loss", portfolioSettingsHandler.SetDailyLossLimit)
		protected.PUT("/portfolio/:id/max-open-trades", portfolioSettingsHandler.SetMaxOpenTrades)
		protected.GET("/portfolio/daily-loss-status", portfolioSettingsHandler.GetDailyLossStatus)

		protected.POST("/trade", tradeHandler.CreateTrade)
		protected.GET("/trade", tradeHandler.GetTrades)
		protected.GET("/trade/:id", tradeHandler.GetTrade)
		protected.PUT("/trade/:id", tradeHandler.UpdateTrade)
		protected.PUT("/trade/:id/close", tradeHandler.CloseTrade)
		protected.PUT("/trade/:id/partial-close", tradeHandler.PartialClose)
		protected.PUT("/trade/:id/breakeven", tradeHandler.MoveToBreakeven)
		protected.DELETE("/trade/:id", tradeHandler.DeleteTrade)

		protected.GET("/dashboard", dashboardHandler.GetDashboard)
		protected.GET("/dashboard/pnl-history", dashboardHandler.GetPnLHistory)
		protected.GET("/dashboard/daily-pnl", dashboardHandler.GetDailyPnL)
		protected.POST("/calculator", dashboardHandler.Calculate)

		protected.POST("/alerts/check", alertHandler.CheckAlerts)
		protected.GET("/alerts", alertHandler.GetUnread)
		protected.PUT("/alerts/:id/read", alertHandler.MarkRead)

		protected.POST("/checklist", checklistHandler.CreateChecklist)
		protected.GET("/checklist/:tradeId", checklistHandler.GetChecklist)
		protected.PUT("/checklist/:tradeId", checklistHandler.UpdateCheck)
		protected.GET("/checklist/defaults", checklistHandler.GetDefaults)
		protected.PUT("/checklist/defaults", checklistHandler.SetDefaults)

		protected.POST("/badges/check", badgeHandler.CheckBadges)
		protected.GET("/badges", badgeHandler.GetBadges)

		protected.GET("/news", newsHandler.GetNews)
		protected.POST("/news/refresh", newsHandler.RefreshNews)
	}

	_ = portfolioRepo

	return r
}
