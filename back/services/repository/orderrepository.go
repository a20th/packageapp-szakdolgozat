package repository

import (
	"back-go/services/models"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	Db *gorm.DB
}

func (o OrderRepository) Delete(s string) error {
	ctx := context.Background()
	_, err := gorm.G[models.Order](o.Db).Where("order_id = ?", s).Delete(ctx)
	return err
}

func (o OrderRepository) Store(order *models.Order) error {
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.Order](o.Db, result, clause.OnConflict{UpdateAll: true}).Create(ctx, order)
	return err
}

func (o OrderRepository) Find(id string) (*models.Order, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Order](o.Db).Where("id = ?", id).First(ctx)
	return &user, err
}

func (o OrderRepository) FindFromAccount(accountEmail string) (*[]models.Order, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Order](o.Db).Where("account_email = ?", accountEmail).Preload("Packages", nil).Order("created_at DESC").Find(ctx)
	return &user, err
}
