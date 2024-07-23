package models

type User struct {
    ID            uint   `gorm:"primaryKey"`
    Email         string `gorm:"unique;not null"`
    AccountNumber string `gorm:"unique;not null"`
}

