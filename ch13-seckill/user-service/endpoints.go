package main

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
)

// CalculateEndpoint define endpoint
type UserEndpoints struct {
	UserEndpoint        endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}

var (
	ErrInvalidRequestType = errors.New("invalid username, password")
)

// ArithmeticRequest define request struct
type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ArithmeticResponse define response struct
type UserResponse struct {
	Result bool  `json:"result"`
	Error  error `json:"error"`
}

//  make endpoint
func MakeUserEndpoint(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(UserRequest)

		var (
			username, password string
			res                bool
			calError           error
		)

		username = req.Username
		password = req.Password

		svc.check(username, password)

		return UserResponse{Result: res, Error: calError}, nil
	}
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
