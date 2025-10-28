package repository

import (
	"back-go/services/models"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountRepository struct {
	Db *gorm.DB
}

func (a AccountRepository) Store(account *models.Account) error {
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.Account](a.Db, result, clause.OnConflict{UpdateAll: true}).Create(ctx, account)
	gorm.G[models.Account](a.Db, result)
	return err
}

func (a AccountRepository) Find(id string) (*models.Account, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Account](a.Db).Where("id = ?", id).First(ctx)
	return &user, err
}

func (a AccountRepository) FindByEmail(email string) (*models.Account, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Account](a.Db).Where("email = ?", email).First(ctx)
	return &user, err
}

func (a AccountRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := gorm.G[models.Account](a.Db).Where("id = ?", id).Delete(ctx)
	return err
}
