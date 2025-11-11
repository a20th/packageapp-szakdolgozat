package pricing

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

func MakeCalculatePriceRequest(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(PriceRequest)
		price, err := s.CalculatePrice(req.From, req.To, req.Size)
		if err != nil {
			return nil, err
		}
		return PriceResponse{Price: price}, nil
	}
}

func DecodeCalculatePriceRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request PriceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

func DecodePricingRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var request Pricing
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return nil, err
	}
	return request, nil
}

type PriceRequest struct {
	From string `json:"from"`
	To   string `json:"to"`
	Size int    `json:"size"`
}

type PriceResponse struct {
	Price float64 `json:"price"`
}

type Pricing struct {
	KmPrice   int `json:"kmprice"`
	BasePrice int `json:"baseprice"`
}
