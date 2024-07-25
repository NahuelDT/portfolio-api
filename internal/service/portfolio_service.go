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

	// Get user's positions
	positions, err := s.orderRepo.GetUserPositions(userID)
	if err != nil {
		return nil, err
	}

	portfolio := &models.Portfolio{
		AvailableCash: cash,
		Assets:        make([]models.PortfolioAsset, 0),
	}

	for _, position := range positions {
		instrument, err := s.instrumentRepo.FindByID(position.InstrumentID)
		if err != nil {
			return nil, err
		}

		marketData, err := s.marketDataRepo.GetLatestMarketData(position.InstrumentID)
		if err != nil {
			return nil, err
		}

		totalValue := position.Size * marketData.Close
		returnPercentage := (marketData.Close - position.Price) / position.Price * 100

		asset := models.PortfolioAsset{
			Ticker:     instrument.Ticker,
			Name:       instrument.Name,
			Quantity:   position.Size,
			TotalValue: totalValue,
			Return:     returnPercentage,
		}

		portfolio.Assets = append(portfolio.Assets, asset)
		portfolio.TotalValue += totalValue
	}

	portfolio.TotalValue += cash

	return portfolio, nil
}
