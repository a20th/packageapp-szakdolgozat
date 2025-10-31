package pricing

import (
	"context"
	"encoding/json"
	"errors"
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

type Service interface {
	CalculatePrice(from string, to string, size int) (float64, error)
}

type Pricing struct {
	BasePrice int `json:"base_price"`
	KmPrice   int `json:"km_price"`
}

type service struct {
	apiKey  string
	rate    *rate.Limiter
	pricing Pricing
}

func (s *service) CalculatePrice(from string, to string, size int) (float64, error) {

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

	_, distance, err := geodist.VincentyDistance(fromCoord, toCoord)
	if err != nil {
		return 0, err
	}
	base := float64(s.pricing.BasePrice)
	kmPrice := float64(s.pricing.KmPrice)
	sizeModifier := float64(1)
	if size > 80 {
		sizeModifier = 1.2
	}
	if size > 50 {
		sizeModifier = 1.1
	}

	price := base + math.Round(distance*kmPrice/50)*50*sizeModifier
	price = math.Round(price)
	return price, nil
}

func (s *service) asyncGetLocation(location string, coord chan geodist.Coord, errc chan error) {

	/*val, ok := s.locationCache.Get(location)
	if ok {
		coord <- val.(geodist.Coord)
		errc <- nil
		return
	}*/
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
	defer resp.Body.Close()

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

func CreatePricingService(apiKey string, limit *rate.Limiter, pricing Pricing) Service {
	return &service{
		apiKey:  apiKey,
		rate:    limit,
		pricing: pricing,
	}
}
