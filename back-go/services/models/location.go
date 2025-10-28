package models

type Location struct {
	Name    string `gorm:"not null"`
	Phone   string `gorm:"not null"`
	Email   *string
	Country string `gorm:"not null"`
	ZIP     string `gorm:"not null"`
	City    string `gorm:"not null"`
	Address string `gorm:"not null"`
	Number  string `gorm:"not null"`
	Other   *string
}

func (l Location) ToCalcString() string {
	return l.Number + "," + l.Address + "," + l.City + "," + l.ZIP + "," + l.Country
}
