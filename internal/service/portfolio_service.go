package service

import (
	"fmt"

	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/NahuelDT/portfolio-api/internal/repository"
)

type PortfolioService struct {
	userRepo       repository.UserRepositorer
	orderRepo      repository.OrderRepositorer
	instrumentRepo repository.InstrumentRepositorer
	marketDataRepo repository.MarketDataRepositorer
}

func NewPortfolioService(
	userRepo repository.UserRepositorer,
	orderRepo repository.OrderRepositorer,
	instrumentRepo repository.InstrumentRepositorer,
	marketDataRepo repository.MarketDataRepositorer,
) *PortfolioService {
	return &PortfolioService{
		userRepo:       userRepo,
		orderRepo:      orderRepo,
		instrumentRepo: instrumentRepo,
		marketDataRepo: marketDataRepo,
	}
}

func (s *PortfolioService) GetPortfolio(userID uint) (*models.Portfolio, error) {
	// Check if user exists
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user with ID %d not found", userID)
	}

	// Get user's cash balance
	cash, err := s.orderRepo.GetUserCashBalance(userID)
	if err != nil {
		return nil, err
	}

	// Get user's filled orders
	orders, err := s.orderRepo.GetUserFilledOrders(userID)
	if err != nil {
		return nil, err
	}

	portfolio := &models.Portfolio{
		AvailableCash: cash,
		Assets:        make([]models.PortfolioAsset, 0),
	}

	// Calculate net positions
	positions := make(map[uint]float64)
	for _, order := range orders {
		if order.Side == "BUY" {
			positions[order.InstrumentID] += order.Size
		} else if order.Side == "SELL" {
			positions[order.InstrumentID] -= order.Size
		}
	}

	for instrumentID, netQuantity := range positions {
		if netQuantity > 0 {
			instrument, err := s.instrumentRepo.GetByID(instrumentID)
			if err != nil {
				return nil, err
			}

			marketData, err := s.marketDataRepo.GetLatestMarketData(instrumentID)
			if err != nil {
				return nil, err
			}

			// Calculate average purchase price
			totalCost := 0.0
			totalQuantity := 0.0
			for _, order := range orders {
				if order.InstrumentID == instrumentID && order.Side == "BUY" {
					totalCost += order.Price * order.Size
					totalQuantity += order.Size
				}
			}
			avgPrice := totalCost / totalQuantity

			totalValue := netQuantity * marketData.Close
			returnPercentage := (marketData.Close - avgPrice) / avgPrice * 100

			asset := models.PortfolioAsset{
				Ticker:     instrument.Ticker,
				Name:       instrument.Name,
				Quantity:   netQuantity,
				TotalValue: totalValue,
				Return:     returnPercentage,
			}

			portfolio.Assets = append(portfolio.Assets, asset)
			portfolio.TotalValue += totalValue
		}
	}

	portfolio.TotalValue += cash

	return portfolio, nil
}
