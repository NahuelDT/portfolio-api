package api

import (
	"github.com/NahuelDT/portfolio-api/internal/api/handlers"
	"github.com/NahuelDT/portfolio-api/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	r *gin.Engine,
	portfolioHandler *handlers.PortfolioHandler,
	searchHandler *handlers.SearchHandler,
	orderHandler *handlers.OrderHandler,
) {
	api := r.Group("/api")
	api.Use(middleware.ErrorHandler())

	api.GET("/portfolio/:userID", portfolioHandler.GetPortfolio)
	api.GET("/search", searchHandler.SearchAssets)
	api.POST("/order", orderHandler.PlaceOrder)
	api.POST("/orders/:orderID/cancel", orderHandler.CancelOrder)
}
