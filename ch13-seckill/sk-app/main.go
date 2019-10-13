package main

import (
	"context"
	"flag"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	register "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/setup"
	localconfig "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/config"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/plugins"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/transport"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var (
		servicePort = flag.String("service.port", bootstrap.HttpConfig.Port, "service port")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	fieldKeys := []string{"method"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "aoho",
		Subsystem: "sk-app",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "aoho",
		Subsystem: "sk-app",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	var svc service.Service
	svc = service.UserService{}

	// add logging middleware
	svc = plugins.LoggingMiddleware(localconfig.Logger)(svc)
	svc = plugins.Metrics(requestCount, requestLatency)(svc)

	userPoint := endpoint.MakeUserEndpoint(svc)
	userPoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(userPoint)
	userPoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "user-endpoint")(userPoint)

	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)
	healthEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "health-endpoint")(healthEndpoint)

	endpts := endpoint.UserEndpoints{
		UserEndpoint:        userPoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, localconfig.ZipkinTracer, localconfig.Logger)

	//http server
	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		setup.InitRedis()
		setup.InitEtcd()
		setup.InitServer()
		//启动前执行注册
		register.Register()
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
	register.Deregister()
	fmt.Println(error)
}
