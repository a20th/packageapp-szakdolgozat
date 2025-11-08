package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	OrderID      string  `gorm:"not null;unique"`
	AccountEmail string  `gorm:"not null"`
	Account      Account `gorm:"not null;foreignKey:AccountEmail;references:Email"`
	Name         string  `gorm:"not null"`
	TaxNumber    *string
	ZIPCode      string    `gorm:"not null"`
	City         string    `gorm:"not null"`
	Country      string    `gorm:"not null"`
	Address      string    `gorm:"not null"`
	Number       string    `gorm:"not null"`
	Packages     []Package `gorm:"foreignKey:OrderID;references:OrderID"`
}
