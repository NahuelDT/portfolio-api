package repository

import (
	"github.com/NahuelDT/portfolio-api/internal/models"
	"gorm.io/gorm"
)

type MarketDataRepository struct {
	db *gorm.DB
}

func NewMarketDataRepository(db *gorm.DB) *MarketDataRepository {
	return &MarketDataRepository{db: db}
}

// GetLatestMarketData retrieves the latest market data for a given instrument
func (r *MarketDataRepository) GetLatestMarketData(instrumentID uint) (*models.MarketData, error) {
	var marketData models.MarketData
	result := r.db.Where("instrumentid = ?", instrumentID).
		Order("date DESC").
		First(&marketData)
	return &marketData, result.Error
}
