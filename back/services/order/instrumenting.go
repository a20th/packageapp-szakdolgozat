package order

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

func (mw InstrumentingMiddleware) GetAllOrders() (orders *[]models.Order, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetAllOrders", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.GetAllOrders()
}

func (mw InstrumentingMiddleware) DeleteOrder(id string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "DeleteOrder", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.DeleteOrder(id)
}

func (mw InstrumentingMiddleware) UpdateOrder(order OrderDTO) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "UpdateOrder", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.UpdateOrder(order)
}

func (mw InstrumentingMiddleware) GetAllOrdersForUser(email string) (orders *[]models.Order, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateOrder", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.GetAllOrdersForUser(email)
}

func (mw InstrumentingMiddleware) CreateOrder(order *models.Order) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "CreateOrder", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return mw.Next.CreateOrder(order)
}

func (mw InstrumentingMiddleware) GetOrder(id string) (order *models.Order, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "GetOrder", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	order, err = mw.Next.GetOrder(id)
	return
}
