package main

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
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/plugins"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	localconfig "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/config"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/endpoint"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var (
		servicePort = flag.String("service.port", bootstrap.HttpConfig.Port, "service port")
		grpcAddr    = flag.String("grpc", ":9008", "gRPC listen address.")
	)

	flag.Parse()

	ctx := context.Background()
	errChan := make(chan error)

	ratebucket := rate.NewLimiter(rate.Every(time.Second*1), 100)

	var tokenService service.TokenService
	var tokenGranter service.TokenGranter
	var clientDetailsService *service.ClientDetailsService

	// add logging middleware


	tokenEndpoint := endpoint.MakeTokenEndpoint(tokenGranter)
	tokenEndpoint = plugins.NewTokenBucketLimitterWithBuildIn(ratebucket)(tokenEndpoint)
	tokenEndpoint = kitzipkin.TraceEndpoint(localconfig.ZipkinTracer, "token-endpoint")(tokenEndpoint)
	tokenEndpoint = plugins.ClientAuthorizationMiddleware(clientDetailsService)(tokenEndpoint)

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
		mysql.InitMysql(conf.MysqlConfig.Host, conf.MysqlConfig.Port, conf.MysqlConfig.User, conf.MysqlConfig.Pwd, conf.MysqlConfig.Db)
		//启动前执行注册
		register.Register()
		handler := r
		errChan <- http.ListenAndServe(":"+*servicePort, handler)
	}()
	//grpc server
	go func() {
		fmt.Println("grpc Server start at port" + *grpcAddr)
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errChan <- err
			return
		}
		serverTracer := kitzipkin.GRPCServerTrace(localconfig.ZipkinTracer, kitzipkin.Name("grpc-transport"))
		tr := localconfig.ZipkinTracer
		md := metadata.MD{}
		parentSpan := tr.StartSpan("test")

		b3.InjectGRPC(&md)(parentSpan.Context())

		ctx := metadata.NewIncomingContext(context.Background(), md)
		handler := transport.NewGRPCServer(ctx, endpts, serverTracer)
		gRPCServer := grpc.NewServer()
		pb.RegisterUserServiceServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
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
