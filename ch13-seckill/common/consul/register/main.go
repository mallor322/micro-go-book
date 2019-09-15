package register

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
)



func RegisterService(svc ConsulService, ctx Context) {

	var (
		consulHost  = flag.String("consul.host", "", "consul ip address")
		consulPort  = flag.String("consul.port", "", "consul port")
		serviceHost = flag.String("service.host", "", "service ip address")
		servicePort = flag.String("service.port", "", "service port")
	)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	//创建健康检查的Endpoint
	healthEndpoint := MakeHealthCheckEndpoint(svc)

	//把算术运算Endpoint和健康检查Endpoint封装至ArithmeticEndpoints
	endpts := HealthCheckEndpoints{
		HealthCheckEndpoint: healthEndpoint,
	}

	//创建http.Handler
	// r := MakeHttpHandler(ctx, endpts, zipkinTracer, logger)
	r := MakeHttpHandler(ctx, endpts, logger)

	//创建注册对象
	registar := Register(*consulHost, *consulPort, *serviceHost, *servicePort, logger)

	go func() {
		fmt.Println("Http Server start at port:" + *servicePort)
		//启动前执行注册
		registar.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()
}