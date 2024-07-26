package repository

import (
	"github.com/NahuelDT/portfolio-api/internal/models"
	"gorm.io/gorm"
)

type InstrumentRepository struct {
	db *gorm.DB
}

func NewInstrumentRepository(db *gorm.DB) *InstrumentRepository {
	return &InstrumentRepository{db: db}
}

// GetByID retrieves an instrument by its ID
func (r *InstrumentRepository) GetByID(id uint) (*models.Instrument, error) {
	var instrument models.Instrument
	result := r.db.First(&instrument, id)
	return &instrument, result.Error
}

// Search performs a general search on instruments based on ticker or name
func (r *InstrumentRepository) Search(query string) ([]models.Instrument, error) {
	var instruments []models.Instrument
	result := r.db.Where("ticker LIKE ? OR name LIKE ?", "%"+query+"%", "%"+query+"%").Find(&instruments)
	return instruments, result.Error
}
