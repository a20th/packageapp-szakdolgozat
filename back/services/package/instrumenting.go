package _package

import (
	"back-go/services/models"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           Service
}

func (mw InstrumentingMiddleware) GetPackageStatus(id string) (status []models.Status, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetPackageStatus", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.GetPackageStatus(id)
}

func (mw InstrumentingMiddleware) AddPackageStatus(id string, status string, description string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "AddPackageStatus", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.AddPackageStatus(id, status, description)
}
