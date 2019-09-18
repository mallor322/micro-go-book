package main

import (
	"github.com/go-kit/kit/log"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/service"
	"time"
)

// loggingMiddleware Make a new type
// that contains Service interface and logger instance
type loggingMiddleware struct {
	service.Service
	logger log.Logger
}

// LoggingMiddleware make logging middleware
func LoggingMiddleware(logger log.Logger) service.ServiceMiddleware {
	return func(next service.Service) service.Service {
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

	ret = mw.Service.Check(a, b)
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
