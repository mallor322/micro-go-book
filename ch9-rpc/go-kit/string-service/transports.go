package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/client"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

func (s *grpcServer) Check(ctx context.Context, r *pb.UserRequest) (*pb.UserResponse, error) {
	_, resp, err := s.check.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.UserResponse), nil
}

func NewGRPCServer(ctx context.Context, endpoints endpts.UserEndpoints, serverTracer grpc.ServerOption) pb.UserServiceServer {
	return &grpcServer{
		check: grpc.NewServer(
			endpoints.UserEndpoint,
			client.DecodeGRPCUserRequest,
			client.EncodeGRPCUserResponse,
			serverTracer,
		),
	}
}

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints StringEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
	}

	r.Methods("POST").Path("/op/{type}/{a}/{b}").Handler(kithttp.NewServer(
		endpoints.StringEndpoint,
		decodeStringRequest,
		encodeStringResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeStringResponse,
		options...,
	))

	return r
}

// decodeStringRequest decode request params to struct
func decodeStringRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	requestType, ok := vars["type"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pa, ok := vars["a"]
	if !ok {
		return nil, ErrorBadRequest
	}

	pb, ok := vars["b"]
	if !ok {
		return nil, ErrorBadRequest
	}

	return StringRequest{
		RequestType: requestType,
		A:           pa,
		B:           pb,
	}, nil
}

// encodeStringResponse encode response to return
func encodeStringResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return HealthRequest{}, nil
}
