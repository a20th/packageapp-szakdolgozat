package _package

import "back-go/services/models"

type Service interface {
	GetPackageStatus(id string) ([]models.Status, error)
	AddPackageStatus(id string, status string, description string) error
}

type service struct {
	Repo Repository
}

func (s service) GetPackageStatus(id string) ([]models.Status, error) {
	pack, err := s.Repo.Find(id)
	if err != nil {
		return nil, err
	}
	return pack.Statuses, err
}

func (s service) AddPackageStatus(id string, status string, description string) error {
	pack, err := s.Repo.Find(id)
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
