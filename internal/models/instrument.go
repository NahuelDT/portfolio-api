package models

type Instrument struct {
	ID     uint   `gorm:"primaryKey;column:id"`
	Ticker string `gorm:"unique;not null:column:ticker"`
	Name   string `gorm:"not null;column:name"`
	Type   string `gorm:"not null;column:type"`
}
