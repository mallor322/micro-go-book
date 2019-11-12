package transport

import (
	"context"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/client"
	endpts "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
)

type grpcServer struct {
	checkTokenServer grpc.Handler
}

func (s *grpcServer) CheckToken(ctx context.Context, r *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	_, resp, err := s.checkTokenServer.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.CheckTokenResponse), nil
}

func NewGRPCServer(ctx context.Context, endpoints endpts.OAuth2Endpoints, serverTracer grpc.ServerOption) pb.OAuthServiceServer {
	return &grpcServer{
		checkTokenServer: grpc.NewServer(
			endpoints.CheckTokenEndpoint,
			client.DecodeGRPCCheckTokenRequest,
			client.EncodeGRPCCheckTokenResponse,
			serverTracer,
		),
	}
}
