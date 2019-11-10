package client

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
	"reflect"
)


type BaseClient interface {

}

type BaseClientImpt struct {

}


type MyClient interface {
	BaseClient


}
// TODO 2019-11-10 加一层client的定义，区分和service的，在client加一些负载均衡的策略

func NewClient(serviceName, request, reponse, reflect.Kind) MyClient  {

	reflect.Value
	Value.
	// ????
	return MyClient
}

func CheckToken(conn *grpc.ClientConn, clientTracer kitgrpc.ClientOption) service.CheckTokenService {


	//con := NewLoadBalanceConn(serviceName)
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