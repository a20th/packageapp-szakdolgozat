package account

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	stdjwt "github.com/golang-jwt/jwt/v4"
)

func MakeRegisterEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(RegisterRequest)
		err := s.Register(req.Email, req.Password, req.Name, req.PhoneNumber)
		return nil, err
	}
}

func MakeGetEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		account, err := s.Get(claims.Subject)
		return AccountResponse{
			Email:       account.Email,
			Name:        account.Name,
			PhoneNumber: account.PhoneNumber,
		}, err
	}
}

func MakeVerifyEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(VerifyRequest)
		err := s.Verify(req.Email, req.Code)
		return nil, err
	}
}

func MakeUpdateEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		get, err := s.Get(claims.Subject)
		if err != nil {
			return nil, err
		}
		req := request.(UpdateRequest)
		if req.Name != nil {
			get.Name = *req.Name
		}
		if req.Email != nil {
			get.Email = *req.Email
		}
		if req.PhoneNumber != nil {
			get.PhoneNumber = *req.PhoneNumber
		}
		err = s.Update(*get)
		return nil, err
	}
}

func MakeDeleteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		claims := ctx.Value(jwt.JWTClaimsContextKey).(*stdjwt.RegisteredClaims)
		err := s.Delete(claims.Subject)
		return nil, err
	}
}

func DecodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err

	}
	return request, nil
}

func DecodeGetRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request EmailRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeRequestVerifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request RequestVerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeVerifyRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type RequestVerifyRequest struct {
	Lang string `json:"lang"`
}

type VerifyRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type UpdateRequest struct {
	Email       *string `json:"email,omitempty"`
	Name        *string `json:"name,omitempty"`
	PhoneNumber *string `json:"phone_number,omitempty"`
}

type AccountResponse struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
}

type EmailRequest struct {
	Email string `json:"email"`
}
