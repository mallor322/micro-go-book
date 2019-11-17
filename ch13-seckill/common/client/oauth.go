package client

import (
	"fmt"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	localconfig "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
	"time"
)

//type BaseClient interface {
//
//}
//
//type BaseClientImpt struct {
//
//}
//
//
//type MyClient interface {
//	BaseClient
//
//
//}
// TODO 2019-11-10 加一层client的定义，区分和service的，在client加一些负载均衡的策略

//func NewClient(serviceName, request, reponse, reflect.Kind) MyClient  {
//
//	reflect.Value
//	Value.
//	// ????
//	return MyClient
//}

func CheckTokenFunc(tokenValue string) (*OAuth2Details, error) {
	serviceName := "user"
	serviceInstance := discover.DiscoveryService(serviceName)

	grpcAddr := fmt.Sprintf("%s:%d", serviceInstance.Host, serviceInstance.Port-1)
	tr := localconfig.ZipkinTracer
	//parentSpan := tr.StartSpan("test")
	//ctx := zipkin.NewContext(context.Background(), parentSpan)
	//todo context for trace info
	clientTracer := kitzipkin.GRPCClientTrace(tr, kitzipkin.Name("grpc-transport"))
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	svr := CheckToken(conn, clientTracer)
	res, err := svr.GetOAuth2DetailsByAccessToken(tokenValue)
	if err != nil {
		fmt.Println("Check error", err.Error())
	}
	return res, err
}

func CheckToken(conn *grpc.ClientConn, clientTracer kitgrpc.ClientOption) service.CheckTokenService {

	var ep = grpctransport.NewClient(conn,
		"pb.OAuthService",
		"CheckToken",
		EncodeGRPCCheckTokenRequest,
		DecodeGRPCCheckTokenResponse,
		pb.CheckTokenResponse{},
		clientTracer,
	).Endpoint()

	return &endpoint.OAuth2Endpoints{
		CheckTokenEndpoint: ep,
	}
}
