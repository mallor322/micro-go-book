package endpoints

import (
	"SecKill/sk_layer/service"
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
)

// CalculateEndpoint define endpoint
type ArithmeticEndpoints struct {
	CalculateEndpoint   endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
	AuthEndpoint        endpoint.Endpoint
}

// ArithmeticRequest define request struct
type ArithmeticRequest struct {
	RequestType string `json:"request_type"`
	A           int    `json:"a"`
	B           int    `json:"b"`
}

// ArithmeticResponse define response struct
type ArithmeticResponse struct {
	Result int   `json:"result"`
	Error  error `json:"error"`
}

// MakeArithmeticEndpoint make endpoint
func MakeArithmeticEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(ArithmeticRequest)

		var (
			res, a, b int
			calError  error
		)

		a = req.A
		b = req.B

		res, calError = svc.SecKill(ctx, req.RequestType, a, b)

		return ArithmeticResponse{Result: res, Error: calError}, nil
	}
}


func (ae ArithmeticEndpoints) Calculate(ctx context.Context, reqType string, a, b int) (res int, err error) {
	//ctx := context.Background()
	resp, err := ae.CalculateEndpoint(ctx, ArithmeticRequest{
		RequestType: reqType,
		A:           a,
		B:           b,
	})
	if err != nil {
		return 0, err
	}
	response := resp.(ArithmeticResponse)
	return response.Result, nil
}

