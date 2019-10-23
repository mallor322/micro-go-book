package setup

import (
	"context"
	"flag"
	"fmt"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	conf "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/config"
	register "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/mysql"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-admin/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-admin/plugins"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-admin/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-admin/transport"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/config"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
	"log"
	"net/http"
	"time"
)

//初始化Http服务
func InitServer(host string) {

	var (
		servicePort = flag.String("service.port", bootstrap.HttpConfig.Port, "service port")
	)

	flag.Parse()

	errChan := make(chan error)

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
	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	var (
		activityService service.ActivityService
		productService  service.ProductService
		skAdminService  service.Service
	)
	skAdminService = service.SkAdminService{}
	activityService = service.ActivityService{}
	productService = service.ProductService{}

	// add logging middleware
	skAdminService = plugins.LoggingMiddleware(config.Logger)(skAdminService)
	skAdminService = plugins.Metrics(requestCount, requestLatency)(skAdminService)

	createActivityEnd := endpoint.MakeCreateActivityEndpoint(activityService)
	createActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createActivityEnd)
	createActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-activity")(createActivityEnd)

	GetActivityEnd := endpoint.MakeGetActivityEndpoint(activityService)
	GetActivityEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetActivityEnd)
	GetActivityEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-activity")(GetActivityEnd)

	createProductEnd := endpoint.MakeCreateProductEndpoint(productService)
	createProductEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(createProductEnd)
	createProductEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "create-product")(createProductEnd)

	GetProductEnd := endpoint.MakeGetProductEndpoint(productService)
	GetProductEnd = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(GetProductEnd)
	GetProductEnd = kitzipkin.TraceEndpoint(config.ZipkinTracer, "get-product")(GetProductEnd)

	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(skAdminService)
	healthEndpoint = kitzipkin.TraceEndpoint(config.ZipkinTracer, "health-endpoint")(healthEndpoint)

	endpts := endpoint.SkAdminEndpoints{
		GetActivityEndpoint:    GetActivityEnd,
		CreateActivityEndpoint: createActivityEnd,
		CreateProductEndpoint:  createProductEnd,
		GetProductEndpoint:     GetProductEnd,
		HealthCheckEndpoint:    healthEndpoint,
	}
	ctx := context.Background()
	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, config.ZipkinTracer, config.Logger)

	//http server
	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.User, conf.MysqlConfig.Pwd, conf.MysqlConfig.Db)
		//启动前执行注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()
}
