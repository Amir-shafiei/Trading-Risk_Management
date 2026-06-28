package main

import (
	"Trading-Risk_Management/config"
	"Trading-Risk_Management/http/handler"
	"Trading-Risk_Management/http/router"
	"Trading-Risk_Management/infrastructure/mysql"
	"Trading-Risk_Management/models"
	"Trading-Risk_Management/repository"
	alertService "Trading-Risk_Management/services/alert"
	authService "Trading-Risk_Management/services/auth"
	badgeService "Trading-Risk_Management/services/badge"
	checklistService "Trading-Risk_Management/services/checklist"
	dashboardService "Trading-Risk_Management/services/dashboard"
	newsService "Trading-Risk_Management/services/news"
	portfolioService "Trading-Risk_Management/services/portfolio"
	tradeService "Trading-Risk_Management/services/trade"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}

	db, err := mysql.Connect(&cfg)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Portfolio{},
		&models.Trade{},
		&models.RefreshToken{},
		&models.Badge{},
		&models.PreTradeChecklist{},
		&models.ChecklistDefaults{},
		&models.RiskAlert{},
		&models.DailyStreak{},
	)
	if err != nil {
		log.Fatal("failed to migrate: ", err)
	}

	userRepo := repository.NewUserRepo(db)
	portfolioRepo := repository.NewPtRepo(db)
	tradeRepo := repository.NewTradeRepo(db)
	refreshTokenRepo := repository.NewRefreshTokenRepo(db)
	alertRepo := repository.NewAlertRepo(db)
	checklistRepo := repository.NewChecklistRepo(db)
	badgeRepo := repository.NewBadgeRepo(db)

	auth := authService.NewAuthService(userRepo, refreshTokenRepo, &cfg)
	portfolio := portfolioService.NewPortfolioService(portfolioRepo)
	badgeSvc := badgeService.NewBadgeService(badgeRepo, tradeRepo)
	trade := tradeService.NewTradeService(tradeRepo, portfolioRepo, badgeSvc)
	dashboardSvc := dashboardService.NewDashboardService(tradeRepo, portfolioRepo)
	alertSvc := alertService.NewAlertService(tradeRepo, portfolioRepo, alertRepo)
	checklistSvc := checklistService.NewChecklistService(checklistRepo)
	newsSvc := newsService.NewNewsService()

	authHandler := handler.NewAuthHandler(auth)
	portfolioHandler := handler.NewPortfolioHandler(portfolio)
	tradeHandler := handler.NewTradeHandler(trade)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)
	alertHandler := handler.NewAlertHandler(alertSvc)
	checklistHandler := handler.NewChecklistHandler(checklistSvc)
	badgeHandler := handler.NewBadgeHandler(badgeSvc)
	portfolioSettingsHandler := handler.NewPortfolioSettingsHandler(portfolioRepo)
	newsHandler := handler.NewNewsHandler(newsSvc)

	r := router.Setup(
		&cfg,
		authHandler,
		tradeHandler,
		portfolioHandler,
		dashboardHandler,
		alertHandler,
		checklistHandler,
		badgeHandler,
		portfolioSettingsHandler,
		newsHandler,
		portfolioRepo,
	)

	log.Printf("server running on port %s", cfg.SERVERPORT)
	if err := r.Run(":" + cfg.SERVERPORT); err != nil {
		log.Fatal("failed to run server: ", err)
	}
}
