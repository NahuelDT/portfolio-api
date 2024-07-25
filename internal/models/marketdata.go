package models

import (
	"time"
)

type MarketData struct {
	ID            uint      `gorm:"primaryKey"`
	InstrumentID  uint      `gorm:"column:instrumentid"`
	High          float64   `gorm:"column:high"`
	Low           float64   `gorm:"column:low"`
	Open          float64   `gorm:"column:open"`
	Close         float64   `gorm:"column:close"`
	PreviousClose float64   `gorm:"column:previousclose"`
	DateTime      time.Time `gorm:"column:date"`
}

// TableName especifica el nombre de la tabla para GORM
func (MarketData) TableName() string {
	return "marketdata"
}
