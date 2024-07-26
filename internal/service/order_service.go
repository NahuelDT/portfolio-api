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

	switch order.Side {
	case "BUY", "SELL":
		// Validate instrument
		if _, err := s.instrumentRepo.GetByID(order.InstrumentID); err != nil {
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
		if order.Size == 0 {
			if totalAmount > 0 {
				order.Size = math.Floor(totalAmount / order.Price)
				if order.Size == 0 {
					return errors.New("insufficient funds for minimum order size")
				}
			} else {
				return errors.New("no order size or total amount provided")
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
			}
		} else { // SELL
			userPositions, err := s.calculateUserPositions(order.UserID)
			if err != nil {
				return err
			}
			availableAssets := userPositions[order.InstrumentID]
			if availableAssets < order.Size {
				order.Status = "REJECTED"
			}
		}

	case "CASH_IN":
		order.Status = "FILLED"

	case "CASH_OUT":
		availableCash, err := s.orderRepo.GetUserCashBalance(order.UserID)
		if err != nil {
			return err
		}
		if availableCash < order.Size {
			order.Status = "REJECTED"
		} else {
			order.Status = "FILLED"
		}

	default:
		return errors.New("invalid order side")
	}

	err := s.orderRepo.Create(order)
	if err != nil {
		return err
	}

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

func (s *OrderService) calculateUserPositions(userID uint) (map[uint]float64, error) {
	filledOrders, err := s.orderRepo.GetUserFilledOrders(userID)
	if err != nil {
		return nil, err // Retorna nil y el error si GetUserFilledOrders falla
	}

	positions := make(map[uint]float64)
	for _, order := range filledOrders {
		if order.Side == "BUY" {
			positions[order.InstrumentID] += order.Size
		} else if order.Side == "SELL" {
			positions[order.InstrumentID] -= order.Size
		}
	}

	return positions, nil
}
