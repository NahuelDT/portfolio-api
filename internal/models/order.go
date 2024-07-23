package models

import (
    "time"
)

type Order struct {
    ID           uint      `gorm:"primaryKey"`
    InstrumentID uint      `gorm:"not null"`
    UserID       uint      `gorm:"not null"`
    Side         string    `gorm:"not null"`
    Size         float64   `gorm:"not null"`
    Price        float64
    Type         string    `gorm:"not null"`
    Status       string    `gorm:"not null"`
    DateTime     time.Time `gorm:"not null"`
}

