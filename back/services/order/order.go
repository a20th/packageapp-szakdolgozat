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
	DeleteOrder(id string) error
	UpdateOrder(order OrderDTO) error
	GetAllOrdersForUser(email string) (*[]models.Order, error)
	GetAllOrders() (*[]models.Order, error)
}

var NoPackageError = errors.New("no package found")

type service struct {
	Repo    Repository
	Pricing pricing.Service
}

func (s service) GetAllOrders() (*[]models.Order, error) {
	account, err := s.Repo.GetAll()
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s service) DeleteOrder(id string) error {
	err := s.Repo.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func (s service) UpdateOrder(order OrderDTO) (err error) {
	o, err := s.Repo.Find(order.Id)
	if err != nil {
		return err
	}
	o.City = order.City
	o.Name = order.Name
	o.Country = order.Country
	o.Address = order.Address
	o.TaxNumber = order.TaxNumber
	o.Active = order.Active
	o.ZIPCode = order.ZIPCode
	o.Number = order.Number
	err = s.Repo.Store(o)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetAllOrdersForUser(email string) (*[]models.Order, error) {
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
	numberOfPackages := len(order.Packages)
	if numberOfPackages == 0 {
		return NoPackageError
	}
	order.OrderID = "O" + ulid.Make().String()
	errs, _ := errgroup.WithContext(context.Background())
	for i := 0; i < numberOfPackages; i++ {
		p := &(order.Packages[i])
		p.PackageID = "P" + ulid.Make().String()
		errs.Go(func() error {
			return s.asyncCalc(p)
		})
	}
	err := errs.Wait()
	if err != nil {
		return err
	}
	order.Active = true
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
	Delete(string) error
	GetAll() (*[]models.Order, error)
}
