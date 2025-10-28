package models

import (
	"time"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Password          string `gorm:"not null"`
	Email             string `gorm:"unique; not null"`
	Name              string `gorm:"not null"`
	PhoneNumber       string `gorm:"not null"`
	PreferredLanguage string `gorm:"not null"`
	VerificationID    *string
	VerifiedAt        *time.Time
}
