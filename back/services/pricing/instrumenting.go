package pricing

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	CountResult    metrics.Histogram
	Next           Service
}

func (mw InstrumentingMiddleware) CalculatePrice(from string, to string, size int) (price float64, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CalculatePrice", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	price, err = mw.Next.CalculatePrice(from, to, size)
	return
}
