package repository

import (
	"back-go/services/models"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type JWTRepository struct {
	Db *gorm.DB
}

func (r JWTRepository) Store(tokens *models.TokenRecord) error {
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.TokenRecord](r.Db, result, clause.OnConflict{UpdateAll: true}).Create(ctx, tokens)
	return err
}

func (r JWTRepository) Find(id string) (*models.TokenRecord, error) {
	ctx := context.Background()
	user, err := gorm.G[models.TokenRecord](r.Db).Where("access_id = ? OR refresh_id = ?", id, id).First(ctx)
	return &user, err
}

func (r JWTRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := gorm.G[models.TokenRecord](r.Db).Where("access_id = ? OR refresh_id = ?", id, id).Delete(ctx)
	return err
}
