package pricing

import (
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   Service
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
