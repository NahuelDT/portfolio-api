package repository

import (
	"github.com/NahuelDT/portfolio-api/internal/models"
)

type UserRepositorer interface {
	GetByID(id uint) (*models.User, error)
	// Añade otros métodos necesarios
}

type OrderRepositorer interface {
	Create(order *models.Order) error
	GetByID(id uint) (*models.Order, error)
	UpdateStatus(orderID uint, status string) error
	GetUserPositions(userID uint) ([]models.Order, error)
	GetUserCashBalance(userID uint) (float64, error)
}

type InstrumentRepositorer interface {
	FindByID(id uint) (*models.Instrument, error)
	Search(query string) ([]models.Instrument, error)
}

type MarketDataRepositorer interface {
	GetLatestMarketData(instrumentID uint) (*models.MarketData, error)
}
