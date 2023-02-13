package main

import (
	"fmt"
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingNMiddleware struct {
	requestCount 	metrics.Counter
	requestLatency	metrics.Histogram
	countResult		metrics.Histogram
	next			StringService
}

func (mv instrumentingNMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		mv.requestCount.With(lvs...).Add(1)
		mv.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mv.next.Uppercase(s)
	return
}

func (mv instrumentingNMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "count", "error", "false"}
		mv.requestCount.With(lvs...).Add(1)
		mv.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
		mv.countResult.Observe(float64(n))
	}(time.Now())

	n = mv.next.Count(s)
	return
}

func (mv instrumentingNMiddleware) JsonToXlsx(s [][]interface{}) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "jsonToXlsx", "error", fmt.Sprint(err != nil)}
		mv.requestCount.With(lvs...).Add(1)
		mv.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())

	output, err = mv.next.JsonToXlsx(s)
	return
}