package endpoint

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch7-discovery/service"
)

type DiscoveryEndpoints struct {
	SayHelloEndpoint	endpoint.Endpoint
	DiscoveryEndpoint	endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

var (
)

type SayHelloRequest struct {
}

type SayHelloResponse struct {
	Message string `json:"message"`

}

//  make endpoint
func MakeSayHelloEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		message := svc.SayHello()
		return SayHelloResponse{
			Message:message,
		}, nil
	}
}



type DiscoveryRequest struct {
	ServiceName string
}


type DiscoveryResponse struct {
	Instances []interface{} `json:"instances"`
	Error string `json:"error"`
}

func MakeDiscoveryEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(DiscoveryRequest)
		instances, err := svc.DiscoveryService(ctx, req.ServiceName)

		var errString = ""
		if err != nil{
			errString = err.Error()
		}

		return &DiscoveryResponse{
			Instances:instances,
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
