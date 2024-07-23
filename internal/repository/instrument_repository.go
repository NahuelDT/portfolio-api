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

// Create inserts a new instrument into the database
func (r *InstrumentRepository) Create(instrument *models.Instrument) error {
	return r.db.Create(instrument).Error
}

// FindByID retrieves an instrument by its ID
func (r *InstrumentRepository) FindByID(id uint) (*models.Instrument, error) {
	var instrument models.Instrument
	result := r.db.First(&instrument, id)
	return &instrument, result.Error
}

// FindByTicker retrieves an instrument by its ticker symbol
func (r *InstrumentRepository) FindByTicker(ticker string) (*models.Instrument, error) {
	var instrument models.Instrument
	result := r.db.Where("ticker = ?", ticker).First(&instrument)
	return &instrument, result.Error
}

// Search performs a general search on instruments based on ticker or name
func (r *InstrumentRepository) Search(query string) ([]models.Instrument, error) {
	var instruments []models.Instrument
	result := r.db.Where("ticker LIKE ? OR name LIKE ?", "%"+query+"%", "%"+query+"%").Find(&instruments)
	return instruments, result.Error
}

// SearchByName performs a search on instruments based on name
func (r *InstrumentRepository) SearchByName(name string) ([]models.Instrument, error) {
	var instruments []models.Instrument
	result := r.db.Where("name LIKE ?", "%"+name+"%").Find(&instruments)
	return instruments, result.Error
}

// Update updates an existing instrument in the database
func (r *InstrumentRepository) Update(instrument *models.Instrument) error {
	return r.db.Save(instrument).Error
}

// Delete removes an instrument from the database
func (r *InstrumentRepository) Delete(id uint) error {
	return r.db.Delete(&models.Instrument{}, id).Error
}

// ListAll retrieves all instruments from the database
func (r *InstrumentRepository) ListAll() ([]models.Instrument, error) {
	var instruments []models.Instrument
	result := r.db.Find(&instruments)
	return instruments, result.Error
}

// ListByType retrieves all instruments of a specific type
func (r *InstrumentRepository) ListByType(instrumentType string) ([]models.Instrument, error) {
	var instruments []models.Instrument
	result := r.db.Where("type = ?", instrumentType).Find(&instruments)
	return instruments, result.Error
}

// Exists checks if an instrument with the given ticker already exists
func (r *InstrumentRepository) Exists(ticker string) (bool, error) {
	var count int64
	result := r.db.Model(&models.Instrument{}).Where("ticker = ?", ticker).Count(&count)
	return count > 0, result.Error
}
