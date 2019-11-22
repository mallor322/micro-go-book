package client

import (
	"context"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/loadbalance"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
)

type UserClient interface {
	CheckUser(ctx context.Context, request *pb.UserRequest) (*pb.UserResponse, error)
}

type UserClientImpl struct {
	/**
	* 可以配置负载均衡策略，重试、等机制。也可以配置invokeAfter和invokerBefore
	 */
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
}

func (impl *UserClientImpl) CheckUser(ctx context.Context, request *pb.UserRequest) (*pb.UserResponse, error) {
	response := new(pb.UserResponse)
	if err := impl.manager.DecoratorInvoke("/pb.UserService/Check", "user_check", ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}
}

func NewUserClient(serviceName string, lb loadbalance.LoadBalance) (UserClient, error) {
	if serviceName == "" {
		serviceName = "user"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &UserClientImpl{
		manager: &DefaultClientManager{
			serviceName: serviceName,
			loadBalance: lb,
		},
		serviceName: serviceName,
		loadBalance: lb,
	}, nil

}
