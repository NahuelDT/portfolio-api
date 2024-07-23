package service

import (
	"strings"

	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/NahuelDT/portfolio-api/internal/repository"
)

type SearchService struct {
	instrumentRepo repository.InstrumentRepositorer
}

func NewSearchService(instrumentRepo repository.InstrumentRepositorer) *SearchService {
	return &SearchService{
		instrumentRepo: instrumentRepo,
	}
}

type SearchResult struct {
	ID     uint   `json:"id"`
	Ticker string `json:"ticker"`
	Name   string `json:"name"`
	Type   string `json:"type"`
}

func (s *SearchService) SearchAssets(query string) ([]SearchResult, error) {
	// Normalize the query
	query = strings.TrimSpace(strings.ToUpper(query))

	// Search for instruments
	instruments, err := s.instrumentRepo.Search(query)
	if err != nil {
		return nil, err
	}

	// Convert instruments to search results
	results := make([]SearchResult, len(instruments))
	for i, instrument := range instruments {
		results[i] = SearchResult{
			ID:     instrument.ID,
			Ticker: instrument.Ticker,
			Name:   instrument.Name,
			Type:   instrument.Type,
		}
	}

	return results, nil
}

func (s *SearchService) GetAssetDetails(id uint) (*models.Instrument, error) {
	return s.instrumentRepo.FindByID(id)
}

func (s *SearchService) SearchByTicker(ticker string) (*SearchResult, error) {
	instrument, err := s.instrumentRepo.FindByTicker(ticker)
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		ID:     instrument.ID,
		Ticker: instrument.Ticker,
		Name:   instrument.Name,
		Type:   instrument.Type,
	}, nil
}

func (s *SearchService) SearchByName(name string) ([]SearchResult, error) {
	instruments, err := s.instrumentRepo.SearchByName(name)
	if err != nil {
		return nil, err
	}

	results := make([]SearchResult, len(instruments))
	for i, instrument := range instruments {
		results[i] = SearchResult{
			ID:     instrument.ID,
			Ticker: instrument.Ticker,
			Name:   instrument.Name,
			Type:   instrument.Type,
		}
	}

	return results, nil
}
