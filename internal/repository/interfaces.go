package repository

import (
	"time"

	"github.com/NahuelDT/portfolio-api/internal/models"
)

type UserRepositorer interface {
	GetByID(id uint) (*models.User, error)
	// Añade otros métodos necesarios
}

type OrderRepositorer interface {
	Create(order *models.Order) error
	GetByID(id uint) (*models.Order, error)
	GetUserOrders(userID uint) ([]models.Order, error)
	GetUserOrdersByStatus(userID uint, status string) ([]models.Order, error)
	UpdateStatus(orderID uint, status string) error
	GetUserPositions(userID uint) ([]models.Order, error)
	GetUserCashBalance(userID uint) (float64, error)
}

type InstrumentRepositorer interface {
	Create(instrument *models.Instrument) error
	FindByID(id uint) (*models.Instrument, error)
	FindByTicker(ticker string) (*models.Instrument, error)
	Search(query string) ([]models.Instrument, error)
	SearchByName(name string) ([]models.Instrument, error)
	Update(instrument *models.Instrument) error
	Delete(id uint) error
	ListAll() ([]models.Instrument, error)
	ListByType(instrumentType string) ([]models.Instrument, error)
	Exists(ticker string) (bool, error)
}

type MarketDataRepositorer interface {
	GetLatestMarketData(instrumentID uint) (*models.MarketData, error)
	GetMarketDataByDateRange(instrumentID uint, startDate, endDate time.Time) ([]models.MarketData, error)
	CreateMarketData(marketData *models.MarketData) error
	UpdateMarketData(marketData *models.MarketData) error
	DeleteMarketData(id uint) error
	GetDailyReturn(instrumentID uint, date time.Time) (float64, error)
}
