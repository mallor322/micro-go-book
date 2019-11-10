package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch10-resiliency/caculate/service"
)

type CalculateEndpoints struct {
	CalculateEndpoint	endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}


type CalculateRequest struct {
	A int
	B int
}

type CalculateResponse struct {
	Result int `json:"result"`

}

//  make endpoint
func MakeCalculateEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		result := svc.Calculate(req.a, req.b)
		return CalculateResponse{
			Result:result,
		}, nil
	}
}


// HealthRequest 健康检查请求结构
type HealthRequest struct{}

// HealthResponse 健康检查响应结构
type HealthResponse struct {
	Status bool `json:"status"`
}

// MakeHealthCheckEndpoint 创建健康检查Endpoint
func MakeHealthCheckEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		status := svc.HealthCheck()
		return HealthResponse{
			Status:status,
		}, nil
	}
}
