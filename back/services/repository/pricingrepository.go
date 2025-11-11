package repository

import (
	"back-go/services/models"
	"context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PricingRepository struct {
	Db *gorm.DB
}

func (p PricingRepository) Store(pricing models.Pricing) error {
	pricing.ID = 1
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.Pricing](p.Db, result, clause.OnConflict{UpdateAll: true}).Create(ctx, &pricing)
	return err
}

func (p PricingRepository) Get() (models.Pricing, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Pricing](p.Db).Where("id = ?", 1).First(ctx)
	return user, err
}
