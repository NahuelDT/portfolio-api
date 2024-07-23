package repository

import (
	"time"

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
	result := r.db.Where("instrument_id = ?", instrumentID).
		Order("date_time DESC").
		First(&marketData)
	return &marketData, result.Error
}

// GetMarketDataByDateRange retrieves market data for a given instrument within a date range
func (r *MarketDataRepository) GetMarketDataByDateRange(instrumentID uint, startDate, endDate time.Time) ([]models.MarketData, error) {
	var marketDataList []models.MarketData
	result := r.db.Where("instrument_id = ? AND date_time BETWEEN ? AND ?", instrumentID, startDate, endDate).
		Order("date_time ASC").
		Find(&marketDataList)
	return marketDataList, result.Error
}

// CreateMarketData inserts a new market data record
func (r *MarketDataRepository) CreateMarketData(marketData *models.MarketData) error {
	return r.db.Create(marketData).Error
}

// UpdateMarketData updates an existing market data record
func (r *MarketDataRepository) UpdateMarketData(marketData *models.MarketData) error {
	return r.db.Save(marketData).Error
}

// DeleteMarketData deletes a market data record
func (r *MarketDataRepository) DeleteMarketData(id uint) error {
	return r.db.Delete(&models.MarketData{}, id).Error
}

// GetDailyReturn calculates the daily return for a given instrument and date
func (r *MarketDataRepository) GetDailyReturn(instrumentID uint, date time.Time) (float64, error) {
	var marketData models.MarketData
	result := r.db.Where("instrument_id = ? AND date_time <= ?", instrumentID, date).
		Order("date_time DESC").
		First(&marketData)
	if result.Error != nil {
		return 0, result.Error
	}

	if marketData.PreviousClose == 0 {
		return 0, nil
	}

	return (marketData.Close - marketData.PreviousClose) / marketData.PreviousClose, nil
}
