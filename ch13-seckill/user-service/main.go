package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	register "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var (
		consulHost  = flag.String("consul.host", "", "consul ip address")
		consulPort  = flag.String("consul.port", "", "consul port")
		serviceHost = flag.String("service.host", "", "service ip address")
		servicePort = flag.String("service.port", "", "service port")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "aoho",
		Subsystem: "user_service",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)

	var svc Service
	svc = UserService{}

	// add logging middleware
	svc = LoggingMiddleware(logger)(svc)
	svc = Metrics(requestCount, requestLatency)(svc)

	endpoint := MakeUserEndpoint(svc)

	//创建健康检查的Endpoint
	healthEndpoint := MakeHealthCheckEndpoint(svc)

	endpts := UserEndpoints{
		UserEndpoint:        endpoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := MakeHttpHandler(ctx, endpts, logger)

	//创建注册对象
	registar := register.Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		//启动前执行注册
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	registar.Deregister()
	fmt.Println(error)
}
