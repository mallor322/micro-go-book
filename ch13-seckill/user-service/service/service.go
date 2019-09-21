package service

import (
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/model"
	"log"
)

// Service Define a service interface
type Service interface {
	Check(username, password string) bool

	// HealthCheck check service health status
	HealthCheck() bool
}

//UserService implement Service interface
type UserService struct {
}

// Add implement check method
func (s UserService) Check(username, password string) bool {
	userEntity := model.NewUserModel()
	res, err := userEntity.CheckUser(username, password)
	if err != nil {
		log.Printf("UserEntity.CreateUser, err : %v", err)
		return false
	}
	return res
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s UserService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
