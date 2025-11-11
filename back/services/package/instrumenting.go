package _package

import (
	"back-go/services/models"
	"back-go/services/order"
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type InstrumentingMiddleware struct {
	RequestCount   metrics.Counter
	RequestLatency metrics.Histogram
	Next           Service
}

func (mw InstrumentingMiddleware) GetPackage(id string) (p *models.Package, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetPackage", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	p, err = mw.Next.GetPackage(id)
	return
}

func (mw InstrumentingMiddleware) UpdatePackage(p order.PackageDTO) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdatePackage", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Next.UpdatePackage(p)
	return
}

func (mw InstrumentingMiddleware) DeletePackage(id string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeletePackage", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	err = mw.Next.DeletePackage(id)
	return
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
