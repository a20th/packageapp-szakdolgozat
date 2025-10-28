package _package

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type AddPackageStatusRequest struct {
	Id          string  `json:"id"`
	Status      string  `json:"status"`
	Description *string `json:"description;omitempty"`
}

type GetPackageStatusRequest struct {
	Id string `json:"id"`
}

type GetPackageStatusResponse struct {
	Status      string  `json:"status"`
	Description *string `json:"description;omitempty"`
	Date        string  `json:"date"`
}

func MakeAddPackageStatus(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddPackageStatusRequest)
		err := s.AddPackageStatus(req.Id, req.Status, *req.Description)
		if err != nil {
			return nil, err
		}
		return nil, nil
	}
}

func MakeGetPackageStatus(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetPackageStatusRequest)
		status, err := s.GetPackageStatus(req.Id)
		response := make([]GetPackageStatusResponse, len(status))
		if err != nil {
			return nil, err
		}
		for i, m := range status {
			response[i] = GetPackageStatusResponse{Status: m.Status, Description: m.Description, Date: m.CreatedAt.String()}
		}
		return response, nil
	}
}
