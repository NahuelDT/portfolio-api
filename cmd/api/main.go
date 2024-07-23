package main

import (
	"log"

	"github.com/NahuelDT/portfolio-api/internal/api"
	"github.com/NahuelDT/portfolio-api/internal/api/handlers"
	"github.com/NahuelDT/portfolio-api/internal/config"
	"github.com/NahuelDT/portfolio-api/internal/repository"
	"github.com/NahuelDT/portfolio-api/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := config.SetupDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	instrumentRepo := repository.NewInstrumentRepository(db)
	marketDataRepo := repository.NewMarketDataRepository(db)

	portfolioService := service.NewPortfolioService(userRepo, orderRepo, instrumentRepo, marketDataRepo)
	searchService := service.NewSearchService(instrumentRepo)
	orderService := service.NewOrderService(orderRepo, userRepo, instrumentRepo, marketDataRepo)

	portfolioHandler := handlers.NewPortfolioHandler(portfolioService)
	searchHandler := handlers.NewSearchHandler(searchService)
	orderHandler := handlers.NewOrderHandler(orderService)

	r := gin.Default()
	api.SetupRoutes(r, portfolioHandler, searchHandler, orderHandler)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
