package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/model"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/service"
	"net/http"
)

// CalculateEndpoint define endpoint
type OAuth2Endpoints struct {
	TokenEndpoint		endpoint.Endpoint
	CheckTokenEndpoint	endpoint.Endpoint
	HealthCheckEndpoint endpoint.Endpoint
}



var (
	ErrInvalidRequestType = errors.New("invalid username, password")
	ForbiddenRequestType = errors.New("invalid ")
)


type CheckTokenRequest struct {
	Token string `json:"token"`
	Reader *http.Request
}

type CheckTokenResponse struct {

	OAuthDetails *model.OAuth2Details

}


type TokenRequest struct {
	GrantType string `json:"grant_type"`
	Reader *http.Request
}


type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
}

//  make endpoint
func MakeTokenEndpoint(svc service.TokenGranter, clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(TokenRequest)

		var clientDetails *model.ClientDetails
		clientId, clientSecret, ok := req.Reader.BasicAuth(); if ok{

			clientDetails, err := clientService.GetClientDetailByClientId(ctx, clientId)

			if err != nil{
				return nil, errors.New("403")
			}

			if !clientDetails.IsMatch(clientId, clientSecret){
				return nil, errors.New("403")
			}

		}else {
			return nil, errors.New("please provide the clientId and clientSecret in authorization")
		}

		token, err := svc.Grant(ctx, req.GrantType, clientDetails, req.Reader); if err == nil {
			return TokenResponse{
				AccessToken:token,
			}, nil
		}
		return nil, err
	}
}



func MakeCheckTokenEndpoint(svc service.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CheckTokenRequest)

		if tokenDetails, err := svc.GetOAuth2DetailsByAccessToken(req.Token); err == nil{
			return CheckTokenResponse{
				OAuthDetails:tokenDetails,
			}, nil
		}

		return nil, err
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
		return HealthResponse{status}, nil
	}
}
