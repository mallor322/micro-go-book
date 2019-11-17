package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	endpts "github.com/keets2012/Micro-Go-Pracrise/ch10-resiliency/caculate/endpoint"
	"net/http"
	"strconv"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter")
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints endpts.CalculateEndpoints, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	options := []kithttp.ServerOption{
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		kithttp.ServerErrorEncoder(encodeError),
	}
    // 算术相加接口 /calculate
	r.Methods("GET").Path("/calculate").Handler(kithttp.NewServer(
		endpoints.CalculateEndpoint,
		decodeCalculateRequest,
		encodeJsonResponse,
		options...,
	))
	// 健康检查接口 /health
	r.Methods("GET").Path("/health").Handler(kithttp.NewServer(
		endpoints.HealthCheckEndpoint,
		decodeHealthCheckRequest,
		encodeJsonResponse,
		options...,
	))
	return r
}
// decodeCalculateRequest 编码请求参数为 CalculateRequest
func decodeCalculateRequest(_ context.Context, r *http.Request) (interface{}, error) {
	a, _ := strconv.Atoi(r.URL.Query().Get("a"))
	b, _ := strconv.Atoi(r.URL.Query().Get("b"))

	return endpts.CalculateRequest{
		A : a,
		B : b,
	}, nil
}
// decodeHealthCheckRequest 编码请求为 HealthRequest
func decodeHealthCheckRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	return endpts.HealthRequest{}, nil
}
// encodeJsonResponse 解码 respose 结构体为 http JSON 响应
func encodeJsonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}
// 解码业务逻辑中出现的 err 到 http 响应
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

