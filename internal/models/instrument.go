package models

type Instrument struct {
    ID     uint   `gorm:"primaryKey"`
    Ticker string `gorm:"unique;not null"`
    Name   string `gorm:"not null"`
    Type   string `gorm:"not null"`
}

