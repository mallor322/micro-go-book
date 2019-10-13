package service

import (
	"context"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/model"
)

// Service Define a service interface
type ClientDetailsService interface {

	GetClientDetailByClientId(ctx context.Context, clientId string)(*model.ClientDetails, error)

}

type ClientService struct {

}

