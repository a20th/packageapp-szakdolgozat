package account

import (
	"back-go/services/models"
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

func (mw InstrumentingMiddleware) Verify(email string, code string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Verify", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Next.Verify(email, code)
	return err
}

func (mw InstrumentingMiddleware) Register(email string, password string, name string, phoneNumber string, preferredLang string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Register", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Next.Register(email, password, name, phoneNumber, preferredLang)
	return err
}

func (mw InstrumentingMiddleware) Get(id string) (output *models.Account, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Get", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mw.Next.Get(id)
	return
}

func (mw InstrumentingMiddleware) Update(account models.Account) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Update", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Next.Update(account)
	return
}

func (mw InstrumentingMiddleware) Delete(id string) (err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "Delete", "error", fmt.Sprint(err != nil)}
		mw.RequestCount.With(lvs...).Add(1)
		mw.RequestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	err = mw.Next.Delete(id)
	return
}
