package order

import (
	"back-go/services/models"
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	stdjwt "github.com/golang-jwt/jwt/v4"
)

func MakeCreateOrderRequest(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateOrderRequest)
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		order := models.Order{
			AccountEmail: claims.Subject,
			Name:         req.Name,
			TaxNumber:    req.TaxNumber,
			ZIPCode:      req.ZIPCode,
			City:         req.City,
			Country:      req.Country,
			Address:      req.Address,
			Number:       req.Number,
		}
		packages := make([]models.Package, len(req.Packages))
		for i := 0; i < len(req.Packages); i++ {
			from := models.Location{
				Name:    req.Packages[i].FromName,
				Phone:   req.Packages[i].FromPhone,
				Email:   &req.Packages[i].FromEmail,
				Country: req.Packages[i].FromCountry,
				ZIP:     req.Packages[i].FromZIP,
				City:    req.Packages[i].FromCity,
				Address: req.Packages[i].FromAddress,
				Number:  req.Packages[i].FromNumber,
				Other:   &req.Packages[i].FromOther,
			}
			to := models.Location{
				Name:    req.Packages[i].ToName,
				Phone:   req.Packages[i].ToPhone,
				Email:   &req.Packages[i].ToEmail,
				Country: req.Packages[i].ToCountry,
				ZIP:     req.Packages[i].ToZIP,
				City:    req.Packages[i].ToCity,
				Address: req.Packages[i].ToAddress,
				Number:  req.Packages[i].ToNumber,
				Other:   &req.Packages[i].ToOther,
			}
			pack := models.Package{
				Length: req.Packages[i].Length,
				Width:  req.Packages[i].Width,
				Height: req.Packages[i].Height,
				From:   from,
				To:     to,
			}
			packages[i] = pack
		}
		order.Packages = &packages

		err := s.CreateOrder(&order)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func DecodeCreateOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request CreateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func MakeGetAllOrdersRequest(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		orders, err := s.GetAllOrders(claims.Subject)
		if err != nil {
			return nil, err
		}
		var response GetAllOrdersResponse
		response.Orders = make([]OrderDTO, len(*orders))
		for i, order := range *orders {
			response.Orders[i] = ToOrderDTO(order)
		}
		return response, nil
	}
}

type CreateOrderRequest struct {
	AccountID uint
	Name      string
	TaxNumber *string
	ZIPCode   string
	City      string
	Country   string
	Address   string
	Number    string
	Packages  []CreatePackageDTO
}

type GetAllOrdersResponse struct {
	Orders []OrderDTO `json:"orders"`
}

type OrderDTO struct {
	Id        string
	Name      string
	TaxNumber *string
	ZIPCode   string
	City      string
	Country   string
	Address   string
	Number    string
	Packages  []PackageDTO
}

type PackageDTO struct {
	Id          string
	Length      int `gorm:"not null"`
	Width       int `gorm:"not null"`
	Height      int `gorm:"not null"`
	FromName    string
	FromPhone   string
	FromEmail   string
	FromCountry string
	FromZIP     string
	FromCity    string
	FromAddress string
	FromNumber  string
	FromOther   string
	ToName      string
	ToPhone     string
	ToEmail     string
	ToCountry   string
	ToZIP       string
	ToCity      string
	ToAddress   string
	ToNumber    string
	ToOther     string
}

func ToOrderDTO(order models.Order) (dto OrderDTO) {
	dto.Id = order.OrderID
	dto.Name = order.Name
	dto.TaxNumber = order.TaxNumber
	dto.ZIPCode = order.ZIPCode
	dto.City = order.City
	dto.Country = order.Country
	dto.Address = order.Address
	dto.Number = order.Number

	dto.Packages = make([]PackageDTO, len(*order.Packages))
	for i, p := range *order.Packages {
		dto.Packages[i] = ToPackageDTO(p)
	}
	return
}

func ToPackageDTO(p models.Package) (dto PackageDTO) {
	dto.Id = p.PackageID
	dto.Length = p.Length
	dto.Width = p.Width
	dto.Height = p.Height
	dto.FromName = p.From.Name
	dto.FromPhone = p.From.Phone
	if p.From.Email != nil {
		dto.FromEmail = *p.From.Email
	}
	dto.FromCountry = p.From.Country
	dto.FromZIP = p.From.ZIP
	dto.FromCity = p.From.City
	dto.FromAddress = p.From.Address
	dto.FromNumber = p.From.Number
	if p.From.Other != nil {
		dto.FromOther = *p.From.Other
	}
	dto.ToName = p.To.Name
	dto.ToPhone = p.To.Phone
	if p.To.Email != nil {
		dto.ToEmail = *p.To.Email
	}
	dto.ToCountry = p.To.Country
	dto.ToZIP = p.To.ZIP
	dto.ToCity = p.To.City
	dto.ToAddress = p.To.Address
	dto.ToNumber = p.To.Number
	if p.To.Other != nil {
		dto.ToOther = *p.To.Other
	}
	return
}

type CreatePackageDTO struct {
	Length      int `gorm:"not null"`
	Width       int `gorm:"not null"`
	Height      int `gorm:"not null"`
	FromName    string
	FromPhone   string
	FromEmail   string
	FromCountry string
	FromZIP     string
	FromCity    string
	FromAddress string
	FromNumber  string
	FromOther   string
	ToName      string
	ToPhone     string
	ToEmail     string
	ToCountry   string
	ToZIP       string
	ToCity      string
	ToAddress   string
	ToNumber    string
	ToOther     string
}
