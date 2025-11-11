package pricing

import (
	"back-go/services/models"
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   Service
}

func (mw LoggingMiddleware) SetPricing(pricing models.Pricing) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "SetPricing",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.SetPricing(pricing)
	return
}

func (mw LoggingMiddleware) GetPricing() (pricing models.Pricing, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetPricing",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	pricing, err = mw.Next.GetPricing()
	return
}

func (mw LoggingMiddleware) CalculatePrice(from string, to string, size int) (price float64, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "calculatePrice",
			"from", from,
			"to", to,
			"output", price,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	price, err = mw.Next.CalculatePrice(from, to, size)
	return
}
