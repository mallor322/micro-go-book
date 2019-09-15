package register

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"strings"
)

// CalculateEndpoint define endpoint
type HealthCheckEndpoints struct {
	HealthCheckEndpoint endpoint.Endpoint
}

// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{status}, nil
	}
}
