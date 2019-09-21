package client

import (
	"context"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/endpoint"
	"github.com/pkg/errors"
)

func EncodeGRPCUserRequest(_ context.Context, r interface{}) (interface{}, error) {
	req := r.(endpoint.UserRequest)
	return &pb.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func DecodeGRPCUserRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.UserRequest)
	return endpoint.UserRequest{
		Username: string(req.Username),
		Password: string(req.Password),
	}, nil
}

func EncodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(endpoint.UserResponse)

	if resp.Error != nil {
		return &pb.UserResponse{
			Result: bool(resp.Result),
			Err:    resp.Error.Error(),
		}, nil
	}

	return &pb.UserResponse{
		Result: bool(resp.Result),
		Err:    "",
	}, nil
}

func DecodeGRPCUserResponse(_ context.Context, r interface{}) (interface{}, error) {
	resp := r.(*pb.UserResponse)
	return endpoint.UserResponse{
		Result: bool(resp.Result),
		Error:  errors.New(resp.Err),
	}, nil
}
