package service

import (
	"context"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/model"
)

// Service Define a service interface
type ClientDetailsService interface {

	GetClientDetailByClientId(ctx context.Context, clientId string)(*model.ClientDetails, error)

}

type MysqlClientDetailsService struct {
}

func NewMysqlClientDetailsService() ClientDetailsService {
	return &MysqlClientDetailsService{}
}

func (service *MysqlClientDetailsService)GetClientDetailByClientId(ctx context.Context, clientId string)(*model.ClientDetails, error) {

	return &model.ClientDetails{
		ClientId:                    "clientId",
		ClientSecret:                "clientSecret",
		AccessTokenValiditySeconds:  1800,
		RefreshTokenValiditySeconds: 18000,
		RegisteredRedirectUri:       "http://127.0.0.1",
		AuthorizedGrantTypes:        [] string{"password", "refresh_token"},
	}, nil
}


