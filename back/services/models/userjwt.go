package models

import (
	"gorm.io/gorm"
)

type TokenRecord struct {
	gorm.Model
	AccessID  string `gorm:"not null;unique"`
	RefreshID string `gorm:"not null;unique"`
}
