package models

import (
    "time"
)

type MarketData struct {
    ID            uint      `gorm:"primaryKey"`
    InstrumentID  uint      `gorm:"not null"`
    High          float64   `gorm:"not null"`
    Low           float64   `gorm:"not null"`
    Open          float64   `gorm:"not null"`
    Close         float64   `gorm:"not null"`
    PreviousClose float64   `gorm:"not null"`
    DateTime      time.Time `gorm:"not null"`
}

