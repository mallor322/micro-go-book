package plugins

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-book/ch13-seckill/oauth-service/service"
)

func ClientAuthorizationMiddleware(clientDetailsService service.ClientDetailsService) endpoint.Middleware  {

	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {

			return next(ctx, request)
		}
	}



}

