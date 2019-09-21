package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	endpts "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/endpoint"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpts.UserEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()

	options := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(kithttp.DefaultErrorEncoder),
	}

	r.Methods("POST").Path("/check/valid").Handler(kithttp.NewServer(
		endpoints.UserEndpoint,
		decodeUserRequest,
		encodeUserResponse,
		options...,
	))

	r.Path("/metrics").Handler(promhttp.Handler())

	// create health check handler
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeUserResponse,
		options...,
	))

	return r
}

// decodeUserRequest decode request params to struct
func decodeUserRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var userRequest endpts.UserRequest
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return nil, err
	}
	return userRequest, nil

}

// encodeArithmeticResponse encode response to return
func encodeUserResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

// decodeHealthCheckRequest decode request
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}
