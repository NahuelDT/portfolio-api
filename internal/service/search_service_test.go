package service

import (
	"fmt"
	"testing"

	mocks "github.com/NahuelDT/portfolio-api/internal/mocks/repository"
	"github.com/NahuelDT/portfolio-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchAssets(t *testing.T) {
	// Crear una instancia del mock
	mockRepo := mocks.NewInstrumentRepositorer(t)

	// Crear una instancia del servicio con el mock
	searchService := NewSearchService(mockRepo)

	// Definir el comportamiento esperado del mock
	mockRepo.On("Search", mock.Anything).Return([]models.Instrument{
		{ID: 1, Ticker: "AAPL", Name: "Apple Inc."},
		{ID: 2, Ticker: "GOOGL", Name: "Alphabet Inc."},
	}, nil)

	// Llamar al método que queremos probar
	results, err := searchService.SearchAssets("A")
	fmt.Println(results)

	// Aserciones
	assert.NoError(t, err)
	assert.Len(t, results, 2)
	assert.Equal(t, "AAPL", results[0].Ticker)
	assert.Equal(t, "GOOGL", results[1].Ticker)

	// Verificar que el mock fue llamado como se esperaba
	mockRepo.AssertExpectations(t)
}

func TestGetAssetDetails(t *testing.T) {
	mockRepo := mocks.NewInstrumentRepositorer(t)
	searchService := NewSearchService(mockRepo)

	expectedInstrument := &models.Instrument{ID: 1, Ticker: "AAPL", Name: "Apple Inc."}
	mockRepo.On("FindByID", uint(1)).Return(expectedInstrument, nil)

	asset, err := searchService.GetAssetDetails(1)

	assert.NoError(t, err)
	assert.NotNil(t, asset)
	assert.Equal(t, "AAPL", asset.Ticker)

	mockRepo.AssertExpectations(t)
}

func TestSearchByTicker(t *testing.T) {
	mockRepo := mocks.NewInstrumentRepositorer(t)
	searchService := NewSearchService(mockRepo)

	expectedInstrument := &models.Instrument{ID: 1, Ticker: "AAPL", Name: "Apple Inc."}
	mockRepo.On("FindByTicker", "AAPL").Return(expectedInstrument, nil)

	result, err := searchService.SearchByTicker("AAPL")

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "AAPL", result.Ticker)

	mockRepo.AssertExpectations(t)
}

func TestSearchByName(t *testing.T) {
	mockRepo := mocks.NewInstrumentRepositorer(t)
	searchService := NewSearchService(mockRepo)

	mockRepo.On("SearchByName", "Apple").Return([]models.Instrument{
		{ID: 1, Ticker: "AAPL", Name: "Apple Inc."},
	}, nil)

	results, err := searchService.SearchByName("Apple")

	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "AAPL", results[0].Ticker)

	mockRepo.AssertExpectations(t)
}
