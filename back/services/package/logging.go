package _package

import (
	"back-go/services/models"
	"back-go/services/order"
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   Service
}

func (mw LoggingMiddleware) GetPackage(id string) (p *models.Package, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetPackage",
			"id", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	p, err = mw.Next.GetPackage(id)
	return
}

func (mw LoggingMiddleware) UpdatePackage(p order.PackageDTO) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "UpdatePackage",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.UpdatePackage(p)
	return
}

func (mw LoggingMiddleware) DeletePackage(id string) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "DeletePackage",
			"id", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.DeletePackage(id)
	return
}

func (mw LoggingMiddleware) GetPackageStatus(id string) (status []models.Status, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "GetPackageStatus",
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
			"method", "AddPackageStatus",
			"id", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	return mw.Next.AddPackageStatus(id, status, description)
}
