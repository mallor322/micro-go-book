package handlers

import (
	"context"

	pb "github.com/metaverse/truss/_example"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService() pb.SkAppServiceServer {
	return skappserviceService{}
}

type skappserviceService struct{}

// Seckill implements Service.
func (s skappserviceService) Seckill(ctx context.Context, in *pb.SecRequest) (*pb.SecResponse, error) {
	var resp pb.SecResponse
	resp = pb.SecResponse{
		// ProductId:
		// UserId:
		// Token:
		// TokenTime:
		// Code:
	}
	return &resp, nil
}

// SecInfo implements Service.
func (s skappserviceService) SecInfo(ctx context.Context, in *pb.SecInfoRequest) (*pb.SecInfoResponse, error) {
	var resp pb.SecInfoResponse
	resp = pb.SecInfoResponse{
		// ProductId:
	}
	return &resp, nil
}

// SecInfoList implements Service.
func (s skappserviceService) SecInfoList(ctx context.Context, in *pb.SecInfoListRequest) (*pb.SecInfoListResponse, error) {
	var resp pb.SecInfoListResponse
	resp = pb.SecInfoListResponse{
		// Out:
	}
	return &resp, nil
}
