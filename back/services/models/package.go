package models

import (
	"math"

	"gorm.io/gorm"
)

type Package struct {
	gorm.Model
	PackageID string   `gorm:"not null; unique"`
	Length    int      `gorm:"not null"`
	Width     int      `gorm:"not null"`
	Height    int      `gorm:"not null"`
	Price     int      `gorm:"not null"`
	From      Location `gorm:"embedded;embeddedPrefix:from_"`
	To        Location `gorm:"embedded;embeddedPrefix:to_"`
	OrderID   uint     `gorm:"not null"`
	Statuses  []Status `gorm:"foreignKey:PackageID;references:PackageID"`
}

func (p Package) Size() int {
	return p.Length + p.Width + p.Height - int(math.Min(math.Min(float64(p.Length), float64(p.Width)), float64(p.Height)))
}

type Status struct {
	gorm.Model
	PackageID   string `gorm:"not null"`
	Status      string `gorm:"not null"`
	Description *string
}
