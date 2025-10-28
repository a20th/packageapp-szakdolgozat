package auth

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/endpoint"
	stdjwt "github.com/golang-jwt/jwt/v4"
)

func MakeLoginEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		access, refresh, err := s.Login(req.Email, req.Password)
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
		access, refresh, err := s.Refresh(claims.Subject, claims.ID)
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
		err := s.Logout(claims.ID)
		return nil, err
	}
}

func DecodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeRefreshRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type VerifyRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

type EmailRequest struct {
	Email string `json:"email"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
