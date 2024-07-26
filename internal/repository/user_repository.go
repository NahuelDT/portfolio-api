package repository

import (
	"github.com/NahuelDT/portfolio-api/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByID retrieves an user by its ID
func (r *UserRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, id)
	return &user, result.Error
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}
