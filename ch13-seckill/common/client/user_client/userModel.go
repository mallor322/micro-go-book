package user_client

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
)

type UserEndpoints struct {
	UserEndpoint        endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

func EncodeGRPCUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(pb.UserRequest)
	return &pb.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func DecodeGRPCUserRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UserRequest)
	return pb.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func EncodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(pb.UserResponse)

	if resp.Err != "" {
		return &pb.UserResponse{
			Result: bool(resp.Result),
			Err:    "error",
		}, nil
	}

	return &pb.UserResponse{
		Result: bool(resp.Result),
		Err:    "",
	}, nil
}

func DecodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.UserResponse)
	return pb.UserResponse{
		Result: bool(resp.Result),
		Err:    resp.Err,
	}, nil
}
