package service

import (
	"errors"
	"math"
	"time"

	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/NahuelDT/portfolio-api/internal/repository"
)

type OrderService struct {
	orderRepo      repository.OrderRepositorer
	userRepo       repository.UserRepositorer
	instrumentRepo repository.InstrumentRepositorer
	marketDataRepo repository.MarketDataRepositorer
}

func NewOrderService(
	orderRepo repository.OrderRepositorer,
	userRepo repository.UserRepositorer,
	instrumentRepo repository.InstrumentRepositorer,
	marketDataRepo repository.MarketDataRepositorer,
) *OrderService {
	return &OrderService{
		orderRepo:      orderRepo,
		userRepo:       userRepo,
		instrumentRepo: instrumentRepo,
		marketDataRepo: marketDataRepo,
	}
}

func (s *OrderService) PlaceOrder(order *models.Order, totalAmount float64) error {
	order.DateTime = time.Now()

	// Validate user
	if _, err := s.userRepo.GetByID(order.UserID); err != nil {
		return errors.New("invalid user")
	}

	// Validate instrument
	if _, err := s.instrumentRepo.FindByID(order.InstrumentID); err != nil {
		return errors.New("invalid instrument")
	}

	// Get latest market data
	marketData, err := s.marketDataRepo.GetLatestMarketData(order.InstrumentID)
	if err != nil {
		return errors.New("failed to get market data")
	}

	// Handle MARKET orders
	if order.Type == "MARKET" {
		order.Price = marketData.Close
		order.Status = "FILLED"
	} else if order.Type == "LIMIT" {
		order.Status = "NEW"
	} else {
		return errors.New("invalid order type")
	}

	// Calculate order size if total investment amount is provided
	if order.Size == 0 && totalAmount > 0 {
		order.Size = math.Floor(totalAmount / order.Price)
		if order.Size == 0 {
			return errors.New("insufficient funds for minimum order size")
		}
	}

	// Validate available funds/assets
	if order.Side == "BUY" {
		availableCash, err := s.orderRepo.GetUserCashBalance(order.UserID)
		if err != nil {
			return err
		}
		if availableCash < order.Size*order.Price {
			order.Status = "REJECTED"
			return s.orderRepo.Create(order)
		}
	} else if order.Side == "SELL" {
		userPositions, err := s.orderRepo.GetUserPositions(order.UserID)
		if err != nil {
			return err
		}
		var availableAssets float64
		for _, position := range userPositions {
			if position.InstrumentID == order.InstrumentID {
				availableAssets = position.Size
				break
			}
		}
		if availableAssets < order.Size {
			order.Status = "REJECTED"
			return s.orderRepo.Create(order)
		}
	} else if order.Side == "CASH_IN" || order.Side == "CASH_OUT" {
		// Handle cash transfers
		order.Status = "FILLED"
		if order.Side == "CASH_OUT" {
			availableCash, err := s.orderRepo.GetUserCashBalance(order.UserID)
			if err != nil {
				return err
			}
			if availableCash < order.Size {
				order.Status = "REJECTED"
				return s.orderRepo.Create(order)
			}
		}
	} else {
		return errors.New("invalid order side")
	}

	// Create the order
	err = s.orderRepo.Create(order)
	if err != nil {
		return err
	}

	// Update user positions if order is FILLED
	if order.Status == "FILLED" {
		err = s.updateUserPositions(order)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *OrderService) updateUserPositions(order *models.Order) error {
	// Implementation depends on how you want to handle this
	// For example, you might want to create a new position or update an existing one
	return nil
}

func (s *OrderService) CancelOrder(orderID uint) error {
	order, err := s.orderRepo.GetByID(orderID)
	if err != nil {
		return err
	}

	if order.Status != "NEW" {
		return errors.New("only NEW orders can be cancelled")
	}

	order.Status = "CANCELLED"
	return s.orderRepo.UpdateStatus(orderID, "CANCELLED")
}
