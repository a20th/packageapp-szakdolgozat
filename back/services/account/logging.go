package account

import (
	"back-go/services/models"
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   Service
}

func (mw LoggingMiddleware) Register(email string, password string, name string, phoneNumber string) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "register",
			"input", email,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.Register(email, password, name, phoneNumber)
	return
}

func (mw LoggingMiddleware) Verify(email string, code string) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "verify",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.Verify(email, code)
	return
}

func (mw LoggingMiddleware) Get(id string) (output *models.Account, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "get",
			"input", id,
			"output", output,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	output, err = mw.Next.Get(id)
	return
}

func (mw LoggingMiddleware) Update(account models.Account) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "update",
			"input", account,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.Update(account)
	return
}

func (mw LoggingMiddleware) Delete(id string) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "delete",
			"input", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.Delete(id)
	return

}
