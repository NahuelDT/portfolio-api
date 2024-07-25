package service

import (
	"github.com/NahuelDT/portfolio-api/internal/models"
)

type OrderServicer interface {
	PlaceOrder(order *models.Order, totalAmount float64) error
	CancelOrder(orderID uint) error
	updateUserPositions(order *models.Order) error
}

type PortfolioServicer interface {
	GetPortfolio(userID uint) (*models.Portfolio, error)
}

type SearchServicer interface {
	SearchAssets(query string) ([]SearchResult, error)
}
