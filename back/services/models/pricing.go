package models

import "gorm.io/gorm"

type Pricing struct {
	gorm.Model
	KmPrice   int
	BasePrice int
}
