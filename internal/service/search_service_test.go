package service

import (
	"testing"

	mocks "github.com/NahuelDT/portfolio-api/internal/mocks/repository"
	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchAssets(t *testing.T) {
	mockRepo := mocks.NewInstrumentRepositorer(t)

	searchService := NewSearchService(mockRepo)

	mockRepo.On("Search", mock.Anything).Return([]models.Instrument{
		{ID: 1, Ticker: "AAPL", Name: "Apple Inc."},
		{ID: 2, Ticker: "GOOGL", Name: "Alphabet Inc."},
	}, nil)

	results, err := searchService.SearchAssets("A")

	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "AAPL", results[0].Ticker)
	assert.Equal(t, "GOOGL", results[1].Ticker)

	mockRepo.AssertExpectations(t)
}
