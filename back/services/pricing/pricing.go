package pricing

import (
	"back-go/services/models"
	"context"
	"encoding/json"
	"errors"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/jftuga/geodist"
	"golang.org/x/time/rate"
)

var ErrInvalidLocation = errors.New("invalid location")
var ErrInvalidFrom = errors.New("invalid from location")
var ErrInvalidTo = errors.New("invalid to location")
var ErrInvalidValue = errors.New("invalid value")

type Service interface {
	CalculatePrice(from string, to string, size int) (float64, error)
	SetPricing(pricing models.Pricing) error
	GetPricing() (models.Pricing, error)
}

type service struct {
	apiKey string
	dev    bool
	rate   *rate.Limiter
	repo   Repository
}

type Repository interface {
	Store(pricing models.Pricing) error
	Get() (models.Pricing, error)
}

func (s *service) SetPricing(pricing models.Pricing) error {
	if pricing.KmPrice < 0 || pricing.BasePrice < 0 {
		return ErrInvalidValue
	}
	err := s.repo.Store(pricing)
	return err
}

func (s *service) GetPricing() (models.Pricing, error) {
	pricing, err := s.repo.Get()
	if err != nil {
		return models.Pricing{}, err
	}
	return pricing, nil
}

func (s *service) CalculatePrice(from string, to string, size int) (float64, error) {
	var distance float64
	{
		if s.dev == false {
			fromChan := make(chan geodist.Coord)
			toChan := make(chan geodist.Coord)
			fromError := make(chan error)
			toError := make(chan error)

			go s.asyncGetLocation(from, fromChan, fromError)
			go s.asyncGetLocation(to, toChan, toError)

			fromCoord, toCoord, fromErr, toErr := <-fromChan, <-toChan, <-fromError, <-toError
			if errors.Is(fromErr, ErrInvalidLocation) {
				return 0, ErrInvalidFrom
			}
			if fromErr != nil {
				return 0, fromErr
			}
			if errors.Is(toErr, ErrInvalidLocation) {
				return 0, ErrInvalidTo
			}
			if toErr != nil {
				return 0, toErr
			}

			_, dist, err := geodist.VincentyDistance(fromCoord, toCoord)
			if err != nil {
				return 0, err
			}
			distance = dist
		} else {
			distance = 25
		}
	}

	pricing, err := s.GetPricing()
	if err != nil {
		return 0, err
	}

	base := float64(pricing.BasePrice)
	kmPrice := float64(pricing.KmPrice)
	sizeModifier := float64(1)
	//Nagy csomag
	if size > 80 {
		sizeModifier = 1.2
	}
	//Közepes csomag
	if size > 50 {
		sizeModifier = 1.1
	}

	price := base + math.Round(distance*kmPrice/50)*50*sizeModifier
	price = math.Round(price)
	return price, nil
}

func (s *service) asyncGetLocation(location string, coord chan geodist.Coord, errc chan error) {
	client := http.Client{}
	apiUrl := "https://geocode.maps.co/search?q=?query?&api_key=?apikey?"
	location = url.QueryEscape(location)
	apiUrl = strings.Replace(apiUrl, "?query?", location, -1)
	apiUrl = strings.Replace(apiUrl, "?apikey?", s.apiKey, -1)
	err := s.rate.Wait(context.Background())
	if err != nil {
		coord <- geodist.Coord{}
		errc <- err
		return
	}
	resp, err := client.Get(apiUrl)
	if err != nil {
		coord <- geodist.Coord{}
		errc <- err
		return
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		coord <- geodist.Coord{}
		errc <- errors.New(resp.Status)
		return
	}

	var response GeoAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if len(response) == 0 {
		coord <- geodist.Coord{}
		errc <- ErrInvalidLocation
		return
	}
	if err != nil {
		coord <- geodist.Coord{}
		errc <- err
		return
	}
	lat, err := strconv.ParseFloat(response[0].Lat, 64)
	if err != nil {
		coord <- geodist.Coord{}
		errc <- err
		return
	}
	lon, err := strconv.ParseFloat(response[0].Lon, 64)
	if err != nil {
		coord <- geodist.Coord{}
		errc <- err
		return
	}
	coords := geodist.Coord{
		Lat: lat,
		Lon: lon,
	}

	coord <- coords
	errc <- nil
	return
}

type GeoAPIResponse []struct {
	PlaceID     int      `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmID       int      `json:"osm_id"`
	BoundingBox []string `json:"boundingbox"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	PlaceRank   int      `json:"place_rank"`
	Category    string   `json:"category"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
}

func CreatePricingService(apiKey string, limit *rate.Limiter, repo Repository, dev bool) Service {
	return &service{
		apiKey: apiKey,
		rate:   limit,
		repo:   repo,
		dev:    dev,
	}
}
