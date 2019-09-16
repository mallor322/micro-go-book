package main

import (
	"github.com/go-kit/kit/log"
	"time"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type loggingMiddleware struct {
	Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) ServiceMiddleware {
	return func(next Service) Service {
		return loggingMiddleware{next, logger}
	}
}

func (mw loggingMiddleware) Check(a, b string) (ret bool) {

	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "Check",
			"username", a,
			"pwd", b,
			"result", ret,
			"took", time.Since(begin),
		)
	}(time.Now())

	ret = mw.Service.check(a, b)
	return ret
}

func (mw loggingMiddleware) HealthCheck() (result bool) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"function", "HealthChcek",
			"result", result,
			"took", time.Since(begin),
		)
	}(time.Now())
	result = mw.Service.HealthCheck()
	return
}
