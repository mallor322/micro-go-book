package client

import (
	grpctransport "github.com/go-kit/kit/transport/grpc"
	kitgrpc "github.com/go-kit/kit/transport/grpc"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/service"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
)

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