package models

type User struct {
	ID            uint   `gorm:"primaryKey;column:id"`
	Email         string `gorm:"unique;not null;column:email"`
	AccountNumber string `gorm:"unique;not null;column:accountnumber"`
}
