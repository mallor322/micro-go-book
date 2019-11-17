package user_client

import (
	"github.com/go-kit/kit/circuitbreaker"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
)

var (
	serviceName = "user"
)

func UserCheck(conn *grpc.ClientConn, clientTracer kitgrpc.ClientOption) UserEndpoints {

	var ep = grpctransport.NewClient(conn,
		"pb.UserService",
		"Check",
		EncodeGRPCUserRequest,
		DecodeGRPCUserResponse,
		pb.UserResponse{},
		clientTracer,
	).Endpoint()

	ep = circuitbreaker.Hystrix("user.check")(ep)
	userEp := UserEndpoints{
		UserEndpoint: ep,
	}
	return userEp
}

/*
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
*/
