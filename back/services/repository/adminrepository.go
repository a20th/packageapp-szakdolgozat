package repository

import (
	"back-go/services/models"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AdminRepository struct {
	Db *gorm.DB
}

func (a AdminRepository) GetAll() ([]models.Admin, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Admin](a.Db).Find(ctx)
	return user, err
}

func (a AdminRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := gorm.G[models.Admin](a.Db).Where("username = ?", id).Delete(ctx)
	return err
}

func (a AdminRepository) Store(account *models.Admin) error {
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.Admin](a.Db, result, clause.OnConflict{UpdateAll: true}).Create(ctx, account)
	return err
}

func (a AdminRepository) Find(id string) (*models.Admin, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Admin](a.Db).Where("username = ?", id).First(ctx)
	return &user, err
}
