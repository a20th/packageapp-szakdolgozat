package repository

import (
	"back-go/services/models"
	"context"

	"gorm.io/gorm"
)

type OrderRepository struct {
	Db *gorm.DB
}

func (o OrderRepository) Store(order *models.Order) error {
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.Order](o.Db, result).Create(ctx, order)
	return err
}

func (o OrderRepository) Find(id string) (*models.Order, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Order](o.Db).Where("id = ?", id).First(ctx)
	return &user, err
}

func (o OrderRepository) FindFromAccount(accountEmail string) (*[]models.Order, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Order](o.Db).Where("account_email = ?", accountEmail).Find(ctx)
	return &user, err
}
