package client

import (
	"context"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/loadbalance"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"log"
	"os"
)

var logger = log.New(os.Stderr, "", log.LstdFlags)
var discoveryClient discover.DiscoveryClient = discover.New("114.67.98.210", "8500")

type OAuthClient interface {
	CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
}

func (impl *OAuthClientImpl) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	response := new(pb.CheckTokenResponse)
	if err := impl.manager.DecoratorInvoke("/pb.OAuthService/CheckToken", "token_check", ctx, request, response); err == nil {
		return response, nil
	} else {
		return nil, err
	}
}
func NewOAuthClient(serviceName string, lb loadbalance.LoadBalance) (OAuthClient, error) {
	if serviceName == "" {
		serviceName = "oauth"
	}
	if lb == nil {
		lb = defaultLoadBalance
	}

	return &OAuthClientImpl{
		manager: &DefaultClientManager{
			serviceName: serviceName,
			loadBalance: lb,
		},
		serviceName: serviceName,
		loadBalance: lb,
	}, nil

}
