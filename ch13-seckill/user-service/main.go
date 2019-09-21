package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	register "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/mysql"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/plugins"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/transport"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	var (
		consulHost  = flag.String("consul.host", "106.15.233.99", "consul ip address")
		consulPort  = flag.String("consul.port", "8500", "consul port")
		serviceHost = flag.String("service.host", "localhost", "service ip address")
		servicePort = flag.String("service.port", "9009", "service port")
		mysqlHost   = flag.String("mysql.host", "106.15.233.99", "consul ip address")
		mysqlPort   = flag.String("mysql.port", "3396", "consul port")
		mysqlUser   = flag.String("mysql.user", "root", "service ip address")
		pwdMysql    = flag.String("mysql.pwd", "root_test", "service port")
		dbMysql     = flag.String("mysql.db", "user", "service port")
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

	var svc service.Service
	svc = service.UserService{}

	// add logging middleware
	svc = plugins.LoggingMiddleware(logger)(svc)
	svc = plugins.Metrics(requestCount, requestLatency)(svc)

	userPoint := endpoint.MakeUserEndpoint(svc)

	//创建健康检查的Endpoint
	healthEndpoint := endpoint.MakeHealthCheckEndpoint(svc)

	endpts := endpoint.UserEndpoints{
		UserEndpoint:        userPoint,
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	r := transport.MakeHttpHandler(ctx, endpts, logger)

	//创建注册对象
	registar := register.Register(*consulHost, *consulPort, *serviceHost, *servicePort, "user_service", logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		mysql.InitMysql(*mysqlHost, *mysqlPort, *mysqlUser, *pwdMysql, *dbMysql)
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
