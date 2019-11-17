package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/config"
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

func (oauth2Endpoints *OAuth2Endpoints) GetOAuth2DetailsByAccessToken(tokenValue string) (*model.OAuth2Details, error)  {

	resp, _ := oauth2Endpoints.CheckTokenEndpoint(context.Background(), &CheckTokenRequest{
		Token:tokenValue,
	})
	response := resp.(CheckTokenResponse)
	var err error
	if response.Error != ""{
		err = errors.New(response.Error)
	}
	return response.OAuthDetails, err
}



var (
	ErrInvalidRequest = errors.New("invalid username, password")
	ErrInvalidClientRequest = errors.New("invalid client message")
)


type CheckTokenRequest struct {
	Token string `json:"token"`
	Reader *http.Request
}

type CheckTokenResponse struct {
	OAuthDetails *model.OAuth2Details `json:"o_auth_details"`
	Error string `json:"error"`

}


type TokenRequest struct {
	GrantType string `json:"grant_type"`
	Reader *http.Request
}


type TokenResponse struct {
	AccessToken *model.OAuth2Token `json:"access_token"`
	Error string `json:"error"`
}

//  make endpoint
func MakeTokenEndpoint(svc service.TokenGranter, clientService service.ClientDetailsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*TokenRequest)

		var clientDetails *model.ClientDetails

		if clientId, clientSecret, ok := req.Reader.BasicAuth(); ok{
			clientDetails, err = clientService.GetClientDetailByClientId(ctx, clientId)
			if err != nil{
				conf.Logger.Log("clientId " + clientId  + " is not existed", ErrInvalidRequest)
				return TokenResponse{
					Error:err.Error(),
				}, nil
			}

			if !clientDetails.IsMatch(clientId, clientSecret){
				conf.Logger.Log("clientId and clientSecret not match", ErrInvalidRequest)
				return TokenResponse{
					Error:ErrInvalidClientRequest.Error(),
				},nil
			}

		}else {
			conf.Logger.Log("Error parse clientId and clientSecret in header", ErrInvalidRequest)
			return TokenResponse{
				Error:ErrInvalidClientRequest.Error(),
			},nil
		}

		token, err := svc.Grant(ctx, req.GrantType, clientDetails, req.Reader)

		var errString = ""
		if err != nil{
			errString = err.Error()
		}

		return TokenResponse{
			AccessToken:token,
			Error:errString,
		}, nil
	}
}



func MakeCheckTokenEndpoint(svc service.TokenService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*CheckTokenRequest)
		tokenDetails, err := svc.GetOAuth2DetailsByAccessToken(req.Token)

		var errString = ""
		if err != nil{
			errString = err.Error()
		}

		return CheckTokenResponse{
			OAuthDetails:tokenDetails,
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
