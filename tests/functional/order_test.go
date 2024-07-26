package functional_tests

import (
	"testing"
	"time"

	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/NahuelDT/portfolio-api/internal/repository"
	"github.com/NahuelDT/portfolio-api/internal/service"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTest(t *testing.T) (*gorm.DB, *service.OrderService, *repository.OrderRepository, *repository.UserRepository, *repository.InstrumentRepository, *repository.MarketDataRepository) {
	dsn := "host=jelani.db.elephantsql.com user=gwxuuoxr password=RHT87Wu0WhMrwy1e7da0OFDaEzMDGtIk dbname=gwxuuoxr port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	orderRepo := repository.NewOrderRepository(db)
	userRepo := repository.NewUserRepository(db)
	instrumentRepo := repository.NewInstrumentRepository(db)
	marketDataRepo := repository.NewMarketDataRepository(db)
	orderService := service.NewOrderService(orderRepo, userRepo, instrumentRepo, marketDataRepo)

	return db, orderService, orderRepo, userRepo, instrumentRepo, marketDataRepo
}

func TestMarketBuyOrder(t *testing.T) {
	db, orderService, orderRepo, userRepo, instrumentRepo, marketDataRepo := setupTest(t)

	user := &models.User{Email: "test@example.com", AccountNumber: "TEST123"}
	err := userRepo.Create(user)
	assert.NoError(t, err)

	cashInOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: 66,
		Side:         "CASH_IN",
		Type:         "MARKET",
		Size:         3000,
	}
	err = orderService.PlaceOrder(cashInOrder, 0)
	assert.NoError(t, err)

	instrument := &models.Instrument{Ticker: "AAPL", Name: "Apple Inc.", Type: "STOCK"}
	err = instrumentRepo.Create(instrument)
	assert.NoError(t, err)

	marketData := &models.MarketData{
		InstrumentID: instrument.ID,
		Close:        150.0,
		DateTime:     time.Now(),
	}
	err = marketDataRepo.Create(marketData)
	assert.NoError(t, err)

	buyOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "BUY",
		Type:         "MARKET",
		Size:         10,
	}
	err = orderService.PlaceOrder(buyOrder, 0)
	assert.NoError(t, err)

	createdOrder, err := orderRepo.GetByID(buyOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "FILLED", createdOrder.Status)
	assert.Equal(t, 150.0, createdOrder.Price)

	balance, err := orderRepo.GetUserCashBalance(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1500.0, balance)

	db.Unscoped().Delete(buyOrder)
	db.Unscoped().Delete(marketData)
	db.Unscoped().Delete(instrument)
	db.Unscoped().Delete(cashInOrder)
	db.Unscoped().Delete(user)
}

func TestMarketSellOrder(t *testing.T) {
	db, orderService, orderRepo, userRepo, instrumentRepo, marketDataRepo := setupTest(t)

	user := &models.User{Email: "test@example.com", AccountNumber: "TEST123"}
	err := userRepo.Create(user)
	assert.NoError(t, err)

	cashInOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: 66,
		Side:         "CASH_IN",
		Type:         "MARKET",
		Size:         3000,
	}
	err = orderService.PlaceOrder(cashInOrder, 0)
	assert.NoError(t, err)

	instrument := &models.Instrument{Ticker: "AAPL", Name: "Apple Inc.", Type: "STOCK"}
	err = instrumentRepo.Create(instrument)
	assert.NoError(t, err)

	marketData := &models.MarketData{
		InstrumentID: instrument.ID,
		Close:        150.0,
		DateTime:     time.Now(),
	}
	err = marketDataRepo.Create(marketData)
	assert.NoError(t, err)

	buyOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "BUY",
		Type:         "MARKET",
		Size:         10,
	}
	err = orderService.PlaceOrder(buyOrder, 0)
	assert.NoError(t, err)

	sellOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "SELL",
		Type:         "MARKET",
		Size:         5,
	}
	err = orderService.PlaceOrder(sellOrder, 0)
	assert.NoError(t, err)

	createdOrder, err := orderRepo.GetByID(sellOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "FILLED", createdOrder.Status)
	assert.Equal(t, 150.0, createdOrder.Price)

	balance, err := orderRepo.GetUserCashBalance(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 2250.0, balance)

	db.Unscoped().Delete(sellOrder)
	db.Unscoped().Delete(buyOrder)
	db.Unscoped().Delete(marketData)
	db.Unscoped().Delete(instrument)
	db.Unscoped().Delete(cashInOrder)
	db.Unscoped().Delete(user)
}

func TestLimitBuyOrder(t *testing.T) {
	db, orderService, orderRepo, userRepo, instrumentRepo, marketDataRepo := setupTest(t)

	user := &models.User{Email: "test@example.com", AccountNumber: "TEST123"}
	err := userRepo.Create(user)
	assert.NoError(t, err)

	cashInOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: 66,
		Side:         "CASH_IN",
		Type:         "MARKET",
		Size:         3000,
	}
	err = orderService.PlaceOrder(cashInOrder, 0)
	assert.NoError(t, err)

	instrument := &models.Instrument{Ticker: "AAPL", Name: "Apple Inc.", Type: "STOCK"}
	err = instrumentRepo.Create(instrument)
	assert.NoError(t, err)

	marketData := &models.MarketData{
		InstrumentID: instrument.ID,
		Close:        150.0,
		DateTime:     time.Now(),
	}
	err = marketDataRepo.Create(marketData)
	assert.NoError(t, err)

	limitOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "BUY",
		Type:         "LIMIT",
		Size:         10,
		Price:        145.0,
	}
	err = orderService.PlaceOrder(limitOrder, 0)
	assert.NoError(t, err)

	createdOrder, err := orderRepo.GetByID(limitOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "NEW", createdOrder.Status)

	newMarketData := &models.MarketData{
		InstrumentID: instrument.ID,
		Close:        145.0,
		DateTime:     time.Now(),
	}
	err = marketDataRepo.Create(newMarketData)
	assert.NoError(t, err)

	limitOrderAfterPriceChange, err := orderRepo.GetByID(limitOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "NEW", limitOrderAfterPriceChange.Status)

	balance, err := orderRepo.GetUserCashBalance(user.ID)
	assert.NoError(t, err)
	assert.InDelta(t, 3000.0, balance, 0.01)

	db.Unscoped().Delete(limitOrder)
	db.Unscoped().Delete(newMarketData)
	db.Unscoped().Delete(marketData)
	db.Unscoped().Delete(instrument)
	db.Unscoped().Delete(cashInOrder)
	db.Unscoped().Delete(user)
}

func TestLimitSellOrder(t *testing.T) {
	db, orderService, orderRepo, userRepo, instrumentRepo, marketDataRepo := setupTest(t)

	user := &models.User{Email: "test@example.com", AccountNumber: "TEST123"}
	err := userRepo.Create(user)
	assert.NoError(t, err)

	cashInOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: 66,
		Side:         "CASH_IN",
		Type:         "MARKET",
		Size:         3000,
	}
	err = orderService.PlaceOrder(cashInOrder, 0)
	assert.NoError(t, err)

	instrument := &models.Instrument{Ticker: "AAPL", Name: "Apple Inc.", Type: "STOCK"}
	err = instrumentRepo.Create(instrument)
	assert.NoError(t, err)

	marketData := &models.MarketData{
		InstrumentID: instrument.ID,
		Close:        150.0,
		DateTime:     time.Now(),
	}
	err = marketDataRepo.Create(marketData)
	assert.NoError(t, err)

	buyOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "BUY",
		Type:         "MARKET",
		Size:         10,
	}
	err = orderService.PlaceOrder(buyOrder, 0)
	assert.NoError(t, err)

	balanceAfterBuy, err := orderRepo.GetUserCashBalance(user.ID)
	assert.NoError(t, err)
	assert.InDelta(t, 1500.0, balanceAfterBuy, 0.01)

	limitSellOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "SELL",
		Type:         "LIMIT",
		Size:         5,
		Price:        160.0,
	}
	err = orderService.PlaceOrder(limitSellOrder, 0)
	assert.NoError(t, err)

	createdOrder, err := orderRepo.GetByID(limitSellOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "NEW", createdOrder.Status)

	finalBalance, err := orderRepo.GetUserCashBalance(user.ID)
	assert.NoError(t, err)
	assert.InDelta(t, 1500.0, finalBalance, 0.01)

	db.Unscoped().Delete(limitSellOrder)
	db.Unscoped().Delete(buyOrder)
	db.Unscoped().Delete(marketData)
	db.Unscoped().Delete(instrument)
	db.Unscoped().Delete(cashInOrder)
	db.Unscoped().Delete(user)
}

func TestInsufficientFundsOrder(t *testing.T) {
	db, orderService, orderRepo, userRepo, instrumentRepo, marketDataRepo := setupTest(t)

	user := &models.User{Email: "test@example.com", AccountNumber: "TEST123"}
	err := userRepo.Create(user)
	assert.NoError(t, err)

	cashInOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: 66,
		Side:         "CASH_IN",
		Type:         "MARKET",
		Size:         1000,
	}
	err = orderService.PlaceOrder(cashInOrder, 0)
	assert.NoError(t, err)

	instrument := &models.Instrument{Ticker: "AAPL", Name: "Apple Inc.", Type: "STOCK"}
	err = instrumentRepo.Create(instrument)
	assert.NoError(t, err)

	marketData := &models.MarketData{
		InstrumentID: instrument.ID,
		Close:        150.0,
		DateTime:     time.Now(),
	}
	err = marketDataRepo.Create(marketData)
	assert.NoError(t, err)

	buyOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "BUY",
		Type:         "MARKET",
		Size:         10,
	}
	err = orderService.PlaceOrder(buyOrder, 0)
	assert.NoError(t, err)

	createdOrder, err := orderRepo.GetByID(buyOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "REJECTED", createdOrder.Status)

	balance, err := orderRepo.GetUserCashBalance(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 1000.0, balance)

	db.Unscoped().Delete(buyOrder)
	db.Unscoped().Delete(marketData)
	db.Unscoped().Delete(instrument)
	db.Unscoped().Delete(cashInOrder)
	db.Unscoped().Delete(user)
}

func TestCancelOrder(t *testing.T) {
	db, orderService, orderRepo, userRepo, instrumentRepo, marketDataRepo := setupTest(t)

	user := &models.User{Email: "test@example.com", AccountNumber: "TEST123"}
	err := userRepo.Create(user)
	assert.NoError(t, err)

	cashInOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: 66,
		Side:         "CASH_IN",
		Type:         "MARKET",
		Size:         3000,
	}
	err = orderService.PlaceOrder(cashInOrder, 0)
	assert.NoError(t, err)

	instrument := &models.Instrument{Ticker: "AAPL", Name: "Apple Inc.", Type: "STOCK"}
	err = instrumentRepo.Create(instrument)
	assert.NoError(t, err)

	marketData := &models.MarketData{
		InstrumentID: instrument.ID,
		Close:        150.0,
		DateTime:     time.Now(),
	}
	err = marketDataRepo.Create(marketData)
	assert.NoError(t, err)

	limitOrder := &models.Order{
		UserID:       user.ID,
		InstrumentID: instrument.ID,
		Side:         "BUY",
		Type:         "LIMIT",
		Size:         10,
		Price:        145.0,
	}
	err = orderService.PlaceOrder(limitOrder, 0)
	assert.NoError(t, err)

	createdOrder, err := orderRepo.GetByID(limitOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "NEW", createdOrder.Status)

	err = orderService.CancelOrder(limitOrder.ID)
	assert.NoError(t, err)

	cancelledOrder, err := orderRepo.GetByID(limitOrder.ID)
	assert.NoError(t, err)
	assert.Equal(t, "CANCELLED", cancelledOrder.Status)

	balance, err := orderRepo.GetUserCashBalance(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, 3000.0, balance)

	db.Unscoped().Delete(limitOrder)
	db.Unscoped().Delete(marketData)
	db.Unscoped().Delete(instrument)
	db.Unscoped().Delete(cashInOrder)
	db.Unscoped().Delete(user)
}
