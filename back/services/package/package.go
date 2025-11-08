package _package

import (
	"back-go/services/models"
	"back-go/services/order"
	"database/sql"
	"errors"
)

var ErrNoPackage = errors.New("no package found")

type Service interface {
	GetPackageStatus(id string) ([]models.Status, error)
	AddPackageStatus(id string, status string, description string) error
	GetPackage(id string) (*models.Package, error)
	UpdatePackage(p order.PackageDTO) error
	DeletePackage(id string) error
}

type service struct {
	Repo Repository
}

func (s service) GetPackage(id string) (*models.Package, error) {
	pack, err := s.Repo.Find(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoPackage
	}
	if err != nil {
		return nil, err
	}
	return pack, nil
}

func (s service) UpdatePackage(p order.PackageDTO) error {
	pack, err := s.Repo.Find(p.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNoPackage
	}
	if err != nil {
		return err
	}

	ToPackageModel(p, pack)
	err = s.Repo.Store(pack)
	return err
}

func (s service) DeletePackage(id string) error {
	err := s.Repo.Delete(id)
	return err
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
	Delete(id string) error
}

func CreatePackageService(repo Repository) Service {
	return &service{Repo: repo}
}

func ToPackageModel(dto order.PackageDTO, p *models.Package) {
	p.PackageID = dto.Id
	p.Length = dto.Length
	p.Width = dto.Width
	p.Height = dto.Height

	p.From.Name = dto.FromName
	p.From.Phone = dto.FromPhone
	if dto.FromEmail != "" {
		p.From.Email = &dto.FromEmail
	}
	p.From.Country = dto.FromCountry
	p.From.ZIP = dto.FromZIP
	p.From.City = dto.FromCity
	p.From.Address = dto.FromAddress
	p.From.Number = dto.FromNumber
	if dto.FromOther != "" {
		p.From.Other = &dto.FromOther
	}

	p.To.Name = dto.ToName
	p.To.Phone = dto.ToPhone
	if dto.ToEmail != "" {
		p.To.Email = &dto.ToEmail
	}
	p.To.Country = dto.ToCountry
	p.To.ZIP = dto.ToZIP
	p.To.City = dto.ToCity
	p.To.Address = dto.ToAddress
	p.To.Number = dto.ToNumber
	if dto.ToOther != "" {
		p.To.Other = &dto.ToOther
	}
	return
}
