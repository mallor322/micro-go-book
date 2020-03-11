package service

import (
	"context"
	"errors"
	"github.com/longjoy/micro-go-book/ch13-seckill/oauth-service/model"
)


var (

	ErrClientMessage = errors.New("invalid client")

)

// Service Define a service interface
type ClientDetailsService interface {

	GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string)(*model.ClientDetails, error)

}




type MysqlClientDetailsService struct {
}

func NewMysqlClientDetailsService() ClientDetailsService {
	return &MysqlClientDetailsService{}
}

var defaultClientDetails = &model.ClientDetails{
	ClientId:                    "clientId",
	ClientSecret:                "clientSecret",
	AccessTokenValiditySeconds:  1800,
	RefreshTokenValiditySeconds: 18000,
	RegisteredRedirectUri:       "http://127.0.0.1",
	AuthorizedGrantTypes:        [] string{"password", "refresh_token"},
}

func (service *MysqlClientDetailsService)GetClientDetailByClientId(ctx context.Context, clientId string, clientSecret string)(*model.ClientDetails, error) {

	if(clientId == defaultClientDetails.ClientId && clientSecret == defaultClientDetails.ClientSecret){
		return defaultClientDetails, nil
	}else {
		return nil, ErrClientMessage
	}


}


