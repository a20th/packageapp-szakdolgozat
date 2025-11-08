package _package

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/go-kit/kit/endpoint"
)

type AddPackageStatusRequest struct {
	Id          string `json:"id"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type GetPackageStatusRequest struct {
	Id string `json:"id"`
}

type GetPackageStatusResponse struct {
	Status      string  `json:"status"`
	Description *string `json:"description,omitempty"`
	Date        string  `json:"date"`
}

func DecodeAddPackageStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request AddPackageStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodeGetPackageStatusRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request GetPackageStatusRequest
	id := r.URL.Query().Get("id")
	if id == "" {
		return nil, io.EOF
	}
	request.Id = id
	return request, nil
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
			response[i] = GetPackageStatusResponse{Status: m.Status, Description: m.Description, Date: m.CreatedAt.Format(time.RFC3339)}
		}
		return response, nil
	}
}
