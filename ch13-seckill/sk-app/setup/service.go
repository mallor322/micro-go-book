package setup

import (
	"context"
	"flag"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	register "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pkg/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/config"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/plugins"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/transport"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//初始化Http服务
func InitServer(host string, servicePort string) {

	log.Printf("port is ", servicePort)

	flag.Parse()

	errChan := make(chan error)

	fieldKeys := []string{"method"}

	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "aoho",
		Subsystem: "sk_app",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)

	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "aoho",
		Subsystem: "sk_app",
		Name:      "request_latency",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	var (
		skAppService service.Service
	)
	skAppService = service.SkAppService{}

	// add logging middleware
	skAppService = plugins.SkAppLoggingMiddleware(config.Logger)(skAppService)
	skAppService = plugins.SkAppMetrics(requestCount, requestLatency)(skAppService)

	healthCheckEnd := endpoint.MakeHealthCheckEndpoint(skAppService)
	healthCheckEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(healthCheckEnd)
	healthCheckEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "heath-check")(healthCheckEnd)

	GetSecInfoEnd := endpoint.MakeSecInfoEndpoint(skAppService)
	GetSecInfoEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetSecInfoEnd)
	GetSecInfoEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "sec-info")(GetSecInfoEnd)

	GetSecInfoListEnd := endpoint.MakeSecInfoListEndpoint(skAppService)
	GetSecInfoListEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetSecInfoListEnd)
	GetSecInfoListEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "sec-info-list")(GetSecInfoListEnd)

	SecKillEnd := endpoint.MakeSecKillEndpoint(skAppService)
	SecKillEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(SecKillEnd)
	SecKillEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "sec-kill")(SecKillEnd)

	endpts := endpoint.SkAppEndpoints{
		SecKillEndpoint:        SecKillEnd,
		HeathCheckEndpoint:     healthCheckEnd,
		GetSecInfoEndpoint:     GetSecInfoEnd,
		GetSecInfoListEndpoint: GetSecInfoListEnd,
	}
	ctx := context.Background()
	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.ZipkinTracer, config.Logger)

	//http server
	go func() {
		fmt.Println("Http Server start at port:" + servicePort)
		//启动前执行注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+servicePort, handler)
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
