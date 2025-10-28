package repository

import (
	"back-go/services/models"
	"context"

	"gorm.io/gorm"
)

type PackageRepository struct {
	Db *gorm.DB
}

func (p PackageRepository) Store(m *models.Package) error {
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.Package](p.Db, result).Create(ctx, m)
	return err
}

func (p PackageRepository) Find(id string) (*models.Package, error) {
	ctx := context.Background()
	user, err := gorm.G[models.Package](p.Db).Where("id = ?", id).First(ctx)
	return &user, err
}
