package order

import (
	"back-go/services/models"
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   Service
}

func (mw LoggingMiddleware) GetAllOrders(email string) (orders *[]models.Order, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetAllOrders",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Next.GetAllOrders(email)
}

func (mw LoggingMiddleware) CreateOrder(order *models.Order) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "CreateOrder",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.CreateOrder(order)
	return
}

func (mw LoggingMiddleware) GetOrder(id string) (order *models.Order, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetOrder",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	order, err = mw.Next.GetOrder(id)
	return
}
