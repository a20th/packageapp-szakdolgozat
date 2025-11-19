package repository

import (
	"back-go/services/models"
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type PackageRepository struct {
	Db *gorm.DB
}

func (p PackageRepository) Delete(id string) error {
	ctx := context.Background()
	_, err := gorm.G[models.Package](p.Db).Where("package_id = ?").Delete(ctx)
	return err
}

func (p PackageRepository) Store(m *models.Package) error {
	ctx := context.Background()
	result := gorm.WithResult()
	err := gorm.G[models.Package](p.Db, result, clause.OnConflict{UpdateAll: true}).Create(ctx, m)
	return err
}

func (p PackageRepository) Find(id string) (*models.Package, error) {
	ctx := context.Background()
	pack, err := gorm.G[models.Package](p.Db).
		Where("package_id = ?", id).
		Preload("Statuses", nil).First(ctx)
	return &pack, err
}
