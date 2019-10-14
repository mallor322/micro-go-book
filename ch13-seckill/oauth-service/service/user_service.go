package service

import (
	"context"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/oauth-service/model"
)

// Service Define a service interface
type UserDetailsService interface {

	// Get UserDetails By username
	GetUserDetailByUsername(ctx context.Context, username string)(*model.UserDetails, error)

}

//UserService implement Service interface
type UserService struct {

}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
