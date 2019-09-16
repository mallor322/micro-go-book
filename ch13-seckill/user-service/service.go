package main

// Service Define a service interface
type Service interface {
	check(username, password string) bool

	// HealthCheck check service health status
	HealthCheck() bool
}

//UserService implement Service interface
type UserService struct {
}

// Add implement check method
func (s UserService) check(username, password string) bool {
	return true
}

// HealthCheck implement Service method
// 用于检查服务的健康状态，这里仅仅返回true
func (s UserService) HealthCheck() bool {
	return true
}

// ServiceMiddleware define service middleware
type ServiceMiddleware func(Service) Service
