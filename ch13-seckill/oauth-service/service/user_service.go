package service

import (
	"context"
	"errors"
	"github.com/longjoy/micro-go-book/ch13-seckill/oauth-service/model"
)

var (
	InvalidAuthentication = errors.New("invalid auth")
)
// Service Define a service interface
type UserDetailsService interface {
	// Get UserDetails By username
	GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error)
}

//UserService implement Service interface
type RemoteUserService struct {
}

// TODO not set RemoteUserService
func (service *RemoteUserService) GetUserDetailByUsername(ctx context.Context, username, password string) (*model.UserDetails, error) {
	//res, err := client.Check(username, password)
	//if err != nil || !res {
	//	return nil, InvalidAuthentication
	//} else {

		return &model.UserDetails{
			UserId:      1,
			Username:    username,
			Password:    password,
			Authorities: []string{"Admin", "Super"},
		}, nil
	//}

}

func NewRemoteUserDetailService() *RemoteUserService {
	return &RemoteUserService{
	}
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
