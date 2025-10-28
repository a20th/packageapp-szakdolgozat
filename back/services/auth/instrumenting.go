package auth

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

func (mw InstrumentingMiddleware) Login(email string, password string) (access string, refresh string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Login", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	access, refresh, err = mw.Next.Login(email, password)
	return
}

func (mw InstrumentingMiddleware) Logout(id string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Logout", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Next.Logout(id)
	return err
}

func (mw InstrumentingMiddleware) Refresh(email string, id string) (access string, refresh string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Refresh", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	access, refresh, err = mw.Next.Refresh(email, id)
	return
}
