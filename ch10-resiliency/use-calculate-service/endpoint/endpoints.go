package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch10-resiliency/use-caculate/service"
)

type UseCalculateEndpoints struct {
	UseCalculateEndpoint	endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}


type UseCalculateRequest struct {
	A int
	B int
}

type UseCalculateResponse struct {
	Result int `json:"result"`
	Error string `json:"error"`

}

//  make endpoint
func MakeUseCalculateEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UseCalculateRequest)
		result, err := svc.UseCalculate(req.A, req.B)
		var errString string
		if err != nil{
			errString = err.Error()
		}
		return UseCalculateResponse{
			Result:result,
			Error:errString,
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
