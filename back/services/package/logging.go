package _package

import (
	"back-go/services/models"
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   Service
}

func (mw LoggingMiddleware) GetPackageStatus(id string) (status []models.Status, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "calculatePrice",
			"id", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Next.GetPackageStatus(id)
}

func (mw LoggingMiddleware) AddPackageStatus(id string, status string, description string) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "calculatePrice",
			"id", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Next.AddPackageStatus(id, status, description)
}
