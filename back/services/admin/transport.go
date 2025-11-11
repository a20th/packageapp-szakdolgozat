package admin

import (
	"back-go/services/models"
	"back-go/services/order"
	_package "back-go/services/package"
	"back-go/services/pricing"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	stdjwt "github.com/golang-jwt/jwt/v4"
)

func MakeGetAllOrdersAdminEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		orders, err := s.GetAllOrders()
		if err != nil {
			return nil, err
		}

		orderDTOs := make([]AdminOrderDTO, len(*orders))
		for i, m := range *orders {
			or := AdminOrderDTO{
				OrderDTO:     order.ToOrderDTO(m),
				AccountEmail: m.AccountEmail,
				AccountName:  m.Account.Name,
				AccountPhone: m.Account.PhoneNumber,
			}
			orderDTOs[i] = or
		}

		return orderDTOs, nil
	}
}

func MakeGetPricingAdminEndpoint(s pricing.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		price, err := s.GetPricing()
		if err != nil {
			return nil, err
		}

		resp := pricing.Pricing{
			KmPrice:   price.KmPrice,
			BasePrice: price.BasePrice,
		}
		return resp, nil
	}
}

func MakeSetPricingAdminEndpoint(s pricing.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(pricing.Pricing)
		price := models.Pricing{
			KmPrice:   req.KmPrice,
			BasePrice: req.BasePrice,
		}
		err := s.SetPricing(price)
		return nil, err
	}
}

func MakeCreateAdminEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		if claims.Subject != "admin" {
			return nil, errors.New("only admin user can create admins")
		}
		err := s.Create(req.Username, req.Password)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func MakeDeleteAdminEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		if claims.Subject != "admin" {
			return nil, errors.New("only admin user can delete admins")
		}
		err := s.Delete(req)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func MakeGetAdminsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		if claims.Subject != "admin" {
			return nil, errors.New("only admin user can get admins")
		}
		users, err := s.GetUsers()
		if err != nil {
			return nil, err
		}
		return users, nil
	}
}

func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(Request)
		access, refresh, err := s.Login(req.Username, req.Password)
		if err != nil {
			return nil, err
		}

		return TokenResponse{
			AccessToken:  access,
			RefreshToken: refresh,
		}, err
	}
}

func MakeRefreshEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		access, refresh, err := s.Refresh(claims.ID, claims.Subject)
		if err != nil {
			return nil, err
		}

		return TokenResponse{
			AccessToken:  access,
			RefreshToken: refresh,
		}, err
	}
}

func MakeLogoutEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		err := s.Logout(claims.Subject)
		return nil, err
	}
}

func MakeAdminGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}
}

func MakeAdminGetOrderEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		data, err := s.GetOrder(req)
		if err != nil {
			return nil, err
		}
		return order.ToOrderDTO(*data), err
	}
}

func MakeAdminUpdateOrderEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(order.OrderDTO)
		err := s.UpdateOrder(req)
		return nil, err
	}
}

func MakeAdminDeleteOrderEndpoint(s order.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		err := s.DeleteOrder(req)
		return nil, err
	}
}

func MakeAdminGetPackageEndpoint(s _package.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		data, err := s.GetPackage(req)
		if err != nil {
			return nil, err
		}
		return order.ToPackageDTO(*data), err
	}
}

func MakeAdminUpdatePackageEndpoint(s _package.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(order.PackageDTO)
		err := s.UpdatePackage(req)
		return nil, err
	}
}

func MakeAdminDeletePackageEndpoint(s _package.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(string)
		err := s.DeletePackage(req)
		return nil, err
	}
}

func MakeAdminAddStatusEndpoint(s _package.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(_package.AddPackageStatusRequest)
		err := s.AddPackageStatus(req.Id, req.Status, req.Description)
		return nil, err

	}
}

func DecodeUpdatePackageRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request order.PackageDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeUpdateOrderRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request order.OrderDTO
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeIDHeader(_ context.Context, r *http.Request) (interface{}, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		return nil, io.EOF
	}
	return id, nil
}

type GetAllResponse struct {
	admins []string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type Request struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AdminOrderDTO struct {
	order.OrderDTO
	AccountEmail string
	AccountName  string
	AccountPhone string
}
