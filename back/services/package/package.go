package _package

import (
	"back-go/services/models"
	"database/sql"
	"errors"
)

var ErrNoPackage = errors.New("no package found")

type Service interface {
	GetPackageStatus(id string) ([]models.Status, error)
	AddPackageStatus(id string, status string, description string) error
}

type service struct {
	Repo Repository
}

func (s service) GetPackageStatus(id string) ([]models.Status, error) {
	pack, err := s.Repo.Find(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoPackage
	}
	if err != nil {
		return nil, err
	}
	return pack.Statuses, err
}

func (s service) AddPackageStatus(id string, status string, description string) error {
	pack, err := s.Repo.Find(id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoPackage
	}
	if err != nil {
		return err
	}
	pack.Statuses = append(pack.Statuses, models.Status{PackageID: id, Description: &description, Status: status})
	err = s.Repo.Store(pack)
	return err
}

type Repository interface {
	Store(*models.Package) error
	Find(id string) (*models.Package, error)
}

func CreatePackageService(repo Repository) Service {
	return &service{Repo: repo}
}
