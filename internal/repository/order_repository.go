package repository

import (
	"time"

	"github.com/NahuelDT/portfolio-api/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create inserts a new order into the database
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	result := r.db.First(&order, id)
	return &order, result.Error
}

// UpdateStatus updates the status of an order
func (r *OrderRepository) UpdateStatus(orderID uint, status string) error {
	return r.db.Model(&models.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

// GetUserPositions retrieves the current positions (filled orders) for a user
func (r *OrderRepository) GetUserPositions(userID uint) ([]models.Order, error) {
	var orders []models.Order
	result := r.db.Where("userid = ? AND status = ? AND side = ?", userID, "FILLED", "BUY").Find(&orders)
	return orders, result.Error
}

// GetUserCashBalance calculates the user's cash balance based on their orders
func (r *OrderRepository) GetUserCashBalance(userID uint) (float64, error) {
	var result struct {
		Balance float64
	}

	err := r.db.Model(&models.Order{}).
		Select("COALESCE(SUM(CASE "+
			"WHEN side = 'CASH_IN' THEN size "+
			"WHEN side = 'CASH_OUT' THEN -size "+
			"WHEN side = 'BUY' THEN -size * price "+
			"WHEN side = 'SELL' THEN size * price "+
			"ELSE 0 END), 0) as balance").
		Where("userid = ? AND status = ?", userID, "FILLED").
		Scan(&result).Error

	if err != nil {
		return 0, err
	}

	return result.Balance, nil
}

// GetOrdersInDateRange retrieves orders for a user within a specific date range
func (r *OrderRepository) GetOrdersInDateRange(userID uint, startDate, endDate time.Time) ([]models.Order, error) {
	var orders []models.Order
	result := r.db.Where("userid = ? AND datetime BETWEEN ? AND ?", userID, startDate, endDate).Find(&orders)
	return orders, result.Error
}

// CancelOrder cancels an order by updating its status
func (r *OrderRepository) CancelOrder(orderID uint) error {
	return r.db.Model(&models.Order{}).Where("id = ? AND status = ?", orderID, "NEW").Update("status", "CANCELLED").Error
}

// GetPendingOrders retrieves all pending (NEW) orders
func (r *OrderRepository) GetPendingOrders() ([]models.Order, error) {
	var orders []models.Order
	result := r.db.Where("status = ?", "NEW").Find(&orders)
	return orders, result.Error
}
