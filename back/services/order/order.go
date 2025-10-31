package order

import (
	"back-go/services/models"
	"back-go/services/pricing"
	"context"
	"errors"

	"github.com/oklog/ulid/v2"
	"golang.org/x/sync/errgroup"
)

type Service interface {
	CreateOrder(order *models.Order) error
	GetOrder(id string) (*models.Order, error)
	GetAllOrders(email string) (*[]models.Order, error)
}

var NoPackageError = errors.New("no package found")

type service struct {
	Repo    Repository
	Pricing pricing.Service
}

func (s service) GetAllOrders(email string) (*[]models.Order, error) {
	account, err := s.Repo.FindFromAccount(email)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s service) asyncCalc(p *models.Package) error {
	price, err := s.Pricing.CalculatePrice(p.From.ToCalcString(), p.To.ToCalcString(), p.Size())
	if err != nil {
		return err
	}
	p.Price = int(price)
	p.Statuses = append(p.Statuses, models.Status{Status: "confirmation"})
	return nil
}

func (s service) CreateOrder(order *models.Order) error {
	numberOfPackages := len(*order.Packages)
	if numberOfPackages == 0 {
		return NoPackageError
	}
	order.OrderID = ulid.Make().String()
	errs, _ := errgroup.WithContext(context.Background())
	for i := 0; i < numberOfPackages; i++ {
		p := &(*order.Packages)[i]
		p.PackageID = ulid.Make().String()
		errs.Go(func() error {
			return s.asyncCalc(p)
		})
	}
	err := errs.Wait()
	if err != nil {
		return err
	}
	err = s.Repo.Store(order)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetOrder(id string) (*models.Order, error) {
	return s.Repo.Find(id)
}

func CreateOrderService(repo Repository, serv pricing.Service) Service {
	return &service{Repo: repo, Pricing: serv}
}

type Repository interface {
	Store(order *models.Order) error
	Find(id string) (*models.Order, error)
	FindFromAccount(string) (*[]models.Order, error)
}
