package models

import (
	"time"
)

type Order struct {
	ID           uint      `gorm:"primaryKey"`
	InstrumentID uint      `gorm:"column:instrumentid"`
	UserID       uint      `gorm:"column:userid"`
	Side         string    `gorm:"column:side"`
	Size         float64   `gorm:"column:size"`
	Price        float64   `gorm:"column:price"`
	Type         string    `gorm:"column:type"`
	Status       string    `gorm:"column:status"`
	DateTime     time.Time `gorm:"column:datetime"`
}
