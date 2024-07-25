package service

import (
	"errors"
	"testing"
	"time"

	mocks "github.com/NahuelDT/portfolio-api/internal/mocks/repository"
	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetPortfolio(t *testing.T) {
	mockUserRepo := new(mocks.UserRepositorer)
	mockOrderRepo := new(mocks.OrderRepositorer)
	mockInstrumentRepo := new(mocks.InstrumentRepositorer)
	mockMarketDataRepo := new(mocks.MarketDataRepositorer)

	portfolioService := NewPortfolioService(mockUserRepo, mockOrderRepo, mockInstrumentRepo, mockMarketDataRepo)

	t.Run("Successful portfolio retrieval", func(t *testing.T) {
		userID := uint(1)
		mockUser := &models.User{ID: userID, Email: "test@example.com"}
		mockOrders := []models.Order{
			{ID: 1, InstrumentID: 1, UserID: userID, Side: "BUY", Size: 10, Price: 100, Type: "MARKET", Status: "FILLED"},
			{ID: 2, InstrumentID: 2, UserID: userID, Side: "BUY", Size: 5, Price: 200, Type: "LIMIT", Status: "FILLED"},
		}
		mockInstruments := map[uint]*models.Instrument{
			1: {ID: 1, Ticker: "AAPL", Name: "Apple Inc."},
			2: {ID: 2, Ticker: "GOOGL", Name: "Alphabet Inc."},
		}
		mockMarketData := map[uint]*models.MarketData{
			1: {ID: 1, InstrumentID: 1, Close: 110, DateTime: time.Now()},
			2: {ID: 2, InstrumentID: 2, Close: 220, DateTime: time.Now()},
		}

		mockUserRepo.On("GetByID", userID).Return(mockUser, nil)
		mockOrderRepo.On("GetUserCashBalance", userID).Return(float64(1000), nil)
		mockOrderRepo.On("GetUserPositions", userID).Return(mockOrders, nil)

		for _, order := range mockOrders {
			mockInstrumentRepo.On("FindByID", order.InstrumentID).Return(mockInstruments[order.InstrumentID], nil)
			mockMarketDataRepo.On("GetLatestMarketData", order.InstrumentID).Return(mockMarketData[order.InstrumentID], nil)
		}

		portfolio, err := portfolioService.GetPortfolio(userID)

		assert.NoError(t, err)
		assert.NotNil(t, portfolio)
		assert.Equal(t, float64(1000), portfolio.AvailableCash)
		assert.Equal(t, float64(3200), portfolio.TotalValue) // 1000 (cash) + (10 * 110) + (5 * 220)
		assert.Len(t, portfolio.Assets, 2)

		// Verificar los assets individualmente
		for _, asset := range portfolio.Assets {
			if asset.Ticker == "AAPL" {
				assert.Equal(t, float64(10), asset.Quantity)
				assert.Equal(t, float64(1100), asset.TotalValue)
				assert.Equal(t, float64(10), asset.Return) // (110 - 100) / 100 * 100
			} else if asset.Ticker == "GOOGL" {
				assert.Equal(t, float64(5), asset.Quantity)
				assert.Equal(t, float64(1100), asset.TotalValue)
				assert.Equal(t, float64(10), asset.Return) // (220 - 200) / 200 * 100
			}
		}

		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("User not found", func(t *testing.T) {
		userID := uint(999)
		mockUserRepo.On("GetByID", userID).Return(nil, errors.New("user not found"))

		portfolio, err := portfolioService.GetPortfolio(userID)

		assert.Error(t, err)
		assert.Nil(t, portfolio)
		assert.Contains(t, err.Error(), "user not found")

		mockUserRepo.AssertExpectations(t)
	})

	t.Run("No positions", func(t *testing.T) {
		userID := uint(2)
		mockUser := &models.User{ID: userID, Email: "test2@example.com"}

		mockUserRepo.On("GetByID", userID).Return(mockUser, nil)
		mockOrderRepo.On("GetUserCashBalance", userID).Return(float64(500), nil)
		mockOrderRepo.On("GetUserPositions", userID).Return([]models.Order{}, nil)

		portfolio, err := portfolioService.GetPortfolio(userID)

		assert.NoError(t, err)
		assert.NotNil(t, portfolio)
		assert.Equal(t, float64(500), portfolio.AvailableCash)
		assert.Equal(t, float64(500), portfolio.TotalValue)
		assert.Len(t, portfolio.Assets, 0)

		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

}
