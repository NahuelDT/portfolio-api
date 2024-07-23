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

func (r *UserRepository) GetByID(id uint) (*models.User, error) {
    var user models.User
    result := r.db.First(&user, id)
    return &user, result.Error
}

