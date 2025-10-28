package auth

import (
	"time"

	"github.com/go-kit/kit/log"
)

type LoggingMiddleware struct {
	Logger log.Logger
	Next   Service
}

func (mw LoggingMiddleware) Login(email string, password string) (access string, refresh string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "login",
			"input", email,
			"output", "tokens",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	access, refresh, err = mw.Next.Login(email, password)
	return
}

func (mw LoggingMiddleware) Logout(id string) (err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "logout",
			"input", id,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.Next.Logout(id)
	return
}

func (mw LoggingMiddleware) Refresh(email string, id string) (access string, refresh string, err error) {
	defer func(begin time.Time) {
		_ = mw.Logger.Log(
			"method", "refresh",
			"input", email,
			"output", "tokens",
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	access, refresh, err = mw.Next.Refresh(email, id)
	return
}
