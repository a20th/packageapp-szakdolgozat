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
