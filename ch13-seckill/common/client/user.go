package client

import (
	"context"
	"fmt"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	localconfig "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	endpts "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/service"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"time"
)

var (
	serviceName = "user"
)

func UserCheck(conn *grpc.ClientConn, clientTracer kitgrpc.ClientOption) service.Service {

	var ep = grpctransport.NewClient(conn,
		"pb.UserService",
		"Check",
		EncodeGRPCUserRequest,
		DecodeGRPCUserResponse,
		pb.UserResponse{},
		clientTracer,
	).Endpoint()

	userEp := endpts.UserEndpoints{
		UserEndpoint: ep,
	}
	return userEp
}

func Check(username, password string) (bool, error) {
	serviceInstance := discover.DiscoveryService(serviceName)

	grpcAddr := fmt.Sprintf("%s:%d", serviceInstance.Host, serviceInstance.Port-1)
	tr := localconfig.ZipkinTracer
	parentSpan := tr.StartSpan("test")
	ctx := zipkin.NewContext(context.Background(), parentSpan)

	clientTracer := kitzipkin.GRPCClientTrace(tr, kitzipkin.Name("grpc-transport"))
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	svr := UserCheck(conn, clientTracer)
	result, err := svr.Check(ctx, username, password)
	if err != nil {
		fmt.Println("Check error", err.Error())
	}
	return result, err
}
