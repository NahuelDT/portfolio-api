package service

import (
	"errors"
	"testing"

	mocks "github.com/NahuelDT/portfolio-api/internal/mocks/repository"
	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPlaceOrder(t *testing.T) {
	// setUp function to reinitialize mocks and service before each test
	setUp := func() (*mocks.OrderRepositorer, *mocks.UserRepositorer, *mocks.InstrumentRepositorer, *mocks.MarketDataRepositorer, *OrderService) {
		mockOrderRepo := new(mocks.OrderRepositorer)
		mockUserRepo := new(mocks.UserRepositorer)
		mockInstrumentRepo := new(mocks.InstrumentRepositorer)
		mockMarketDataRepo := new(mocks.MarketDataRepositorer)
		orderService := NewOrderService(mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo)
		return mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService
	}

	t.Run("Place valid MARKET BUY order", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "MARKET",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)
		mockOrderRepo.On("GetUserCashBalance", uint(1)).Return(float64(2000), nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "FILLED", order.Status)
		assert.Equal(t, float64(100), order.Price)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place valid LIMIT BUY order", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "LIMIT",
			Size:         10,
			Price:        90,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)
		mockOrderRepo.On("GetUserCashBalance", uint(1)).Return(float64(2000), nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "NEW", order.Status)
		assert.Equal(t, float64(90), order.Price)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place MARKET BUY order with insufficient funds", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "MARKET",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)
		mockOrderRepo.On("GetUserCashBalance", uint(1)).Return(float64(5), nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "REJECTED", order.Status)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place valid MARKET SELL order", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "SELL",
			Type:         "MARKET",
			Size:         5,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)
		mockOrderRepo.On("GetUserFilledOrders", uint(1)).Return([]models.Order{
			{InstrumentID: 1, Side: "BUY", Size: 10, Status: "FILLED"},
		}, nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "FILLED", order.Status)
		assert.Equal(t, float64(100), order.Price)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place MARKET SELL order with insufficient assets", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "SELL",
			Type:         "MARKET",
			Size:         15,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)
		mockOrderRepo.On("GetUserFilledOrders", uint(1)).Return([]models.Order{
			{InstrumentID: 1, Side: "BUY", Size: 10, Status: "FILLED"},
		}, nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "REJECTED", order.Status)
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place CASH_IN order", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, _, _, orderService := setUp()

		order := &models.Order{
			UserID: 1,
			Side:   "CASH_IN",
			Size:   1000,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "FILLED", order.Status)
		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Place CASH_OUT order with sufficient funds", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, _, _, orderService := setUp()

		order := &models.Order{
			UserID: 1,
			Side:   "CASH_OUT",
			Size:   500,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockOrderRepo.On("GetUserCashBalance", uint(1)).Return(float64(1000), nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "FILLED", order.Status)
		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Place CASH_OUT order with insufficient funds", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, _, _, orderService := setUp()

		order := &models.Order{
			UserID: 1,
			Side:   "CASH_OUT",
			Size:   2000,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockOrderRepo.On("GetUserCashBalance", uint(1)).Return(float64(1000), nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(nil)

		err := orderService.PlaceOrder(order, 0)

		assert.NoError(t, err)
		assert.Equal(t, "REJECTED", order.Status)
		mockUserRepo.AssertExpectations(t)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Place order with invalid user", func(t *testing.T) {
		_, mockUserRepo, _, _, orderService := setUp()

		order := &models.Order{
			UserID:       999,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "MARKET",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(999)).Return(nil, errors.New("user not found"))

		err := orderService.PlaceOrder(order, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid user")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Place order with invalid instrument", func(t *testing.T) {
		_, mockUserRepo, mockInstrumentRepo, _, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 999,
			Side:         "BUY",
			Type:         "MARKET",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(999)).Return(nil, errors.New("instrument not found"))

		err := orderService.PlaceOrder(order, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid instrument")
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
	})

	t.Run("Place order with market data fetch error", func(t *testing.T) {
		_, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "MARKET",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(nil, errors.New("market data fetch error"))

		err := orderService.PlaceOrder(order, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to get market data")
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place order with invalid order type", func(t *testing.T) {
		_, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "INVALID_TYPE",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)

		err := orderService.PlaceOrder(order, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order type")
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place order with invalid order side", func(t *testing.T) {
		_, mockUserRepo, _, _, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "INVALID_SIDE",
			Type:         "MARKET",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)

		err := orderService.PlaceOrder(order, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid order side")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Place order with zero size and no total amount", func(t *testing.T) {
		_, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "MARKET",
			Size:         0,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)

		err := orderService.PlaceOrder(order, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no order size or total amount provided")
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

	t.Run("Place order with create error", func(t *testing.T) {
		mockOrderRepo, mockUserRepo, mockInstrumentRepo, mockMarketDataRepo, orderService := setUp()

		order := &models.Order{
			UserID:       1,
			InstrumentID: 1,
			Side:         "BUY",
			Type:         "MARKET",
			Size:         10,
		}

		mockUserRepo.On("GetByID", uint(1)).Return(&models.User{ID: 1}, nil)
		mockInstrumentRepo.On("GetByID", uint(1)).Return(&models.Instrument{ID: 1}, nil)
		mockMarketDataRepo.On("GetLatestMarketData", uint(1)).Return(&models.MarketData{Close: 100}, nil)
		mockOrderRepo.On("GetUserCashBalance", uint(1)).Return(float64(2000), nil)
		mockOrderRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(errors.New("create error"))

		err := orderService.PlaceOrder(order, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "create error")
		mockOrderRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockInstrumentRepo.AssertExpectations(t)
		mockMarketDataRepo.AssertExpectations(t)
	})

}

func TestCancelOrder(t *testing.T) {
	setUp := func() (*mocks.OrderRepositorer, *OrderService) {
		mockOrderRepo := new(mocks.OrderRepositorer)
		orderService := NewOrderService(mockOrderRepo, nil, nil, nil)
		return mockOrderRepo, orderService
	}

	t.Run("Cancel NEW order", func(t *testing.T) {
		mockOrderRepo, orderService := setUp()

		orderID := uint(1)
		mockOrderRepo.On("GetByID", orderID).Return(&models.Order{ID: orderID, Status: "NEW"}, nil)
		mockOrderRepo.On("UpdateStatus", orderID, "CANCELLED").Return(nil)

		err := orderService.CancelOrder(orderID)

		assert.NoError(t, err)
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Try to cancel FILLED order", func(t *testing.T) {
		mockOrderRepo, orderService := setUp()

		orderID := uint(2)
		mockOrderRepo.On("GetByID", orderID).Return(&models.Order{ID: orderID, Status: "FILLED"}, nil)

		err := orderService.CancelOrder(orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "only NEW orders can be cancelled")
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Try to cancel non-existent order", func(t *testing.T) {
		mockOrderRepo, orderService := setUp()

		orderID := uint(999) // ID de orden que no existe
		mockOrderRepo.On("GetByID", orderID).Return(nil, errors.New("order not found"))

		err := orderService.CancelOrder(orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "order not found")
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Cancel order with update error", func(t *testing.T) {
		mockOrderRepo, orderService := setUp()

		orderID := uint(1)
		mockOrderRepo.On("GetByID", orderID).Return(&models.Order{ID: orderID, Status: "NEW"}, nil)
		mockOrderRepo.On("UpdateStatus", orderID, "CANCELLED").Return(errors.New("update error"))

		err := orderService.CancelOrder(orderID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update error")
		mockOrderRepo.AssertExpectations(t)
	})
}

func TestCalculateUserPositions(t *testing.T) {
	setUp := func() (*mocks.OrderRepositorer, *OrderService) {
		mockOrderRepo := new(mocks.OrderRepositorer)
		orderService := NewOrderService(mockOrderRepo, nil, nil, nil)
		return mockOrderRepo, orderService
	}

	t.Run("Calculate user positions", func(t *testing.T) {
		mockOrderRepo, orderService := setUp()

		userID := uint(1)
		mockOrderRepo.On("GetUserFilledOrders", userID).Return([]models.Order{
			{InstrumentID: 1, Side: "BUY", Size: 10, Status: "FILLED"},
			{InstrumentID: 1, Side: "SELL", Size: 5, Status: "FILLED"},
			{InstrumentID: 2, Side: "BUY", Size: 20, Status: "FILLED"},
		}, nil)

		positions, err := orderService.calculateUserPositions(userID)

		assert.NoError(t, err)
		assert.Equal(t, float64(5), positions[1])
		assert.Equal(t, float64(20), positions[2])
		mockOrderRepo.AssertExpectations(t)
	})

	t.Run("Calculate user positions with error", func(t *testing.T) {
		mockOrderRepo, orderService := setUp()

		userID := uint(1)
		expectedError := errors.New("database error")
		mockOrderRepo.On("GetUserFilledOrders", userID).Return([]models.Order{}, expectedError)

		positions, err := orderService.calculateUserPositions(userID)

		assert.Error(t, err)
		assert.Nil(t, positions)
		assert.Equal(t, expectedError, err)
		mockOrderRepo.AssertExpectations(t)
	})
}
