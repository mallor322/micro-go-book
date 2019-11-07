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
type RemoteUserService struct {

}

func (service *RemoteUserService) GetUserDetailByUsername(ctx context.Context, username string)(*model.UserDetails, error){
	return &model.UserDetails{
		UserId:1,
		Username:username,
		Password:"password",
		Authorities: []string{"Admin", "Super"},
	}, nil
}


func NewRemoteUserDetailService() *RemoteUserService {
	return &RemoteUserService{
	}
}


// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
