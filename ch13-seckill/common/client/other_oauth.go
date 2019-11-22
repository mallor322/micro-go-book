package other_client

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/loadbalance"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

var logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
var defaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalance{}
var discoveryClient discover.DiscoveryClient = discover.New("114.67.98.210", "8500")

type OAuthClient interface {
	CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	/**
	* 可以配置负载均衡策略，重试、等机制。也可以配置invokeAfter和invokerBefore
	 */
	manager     ClientManager
	serviceName string
	loadBalance loadbalance.LoadBalance
}

func (impl *OAuthClientImpl) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	response := new(pb.CheckTokenResponse)
	if err := impl.manager.DecoratorInvoke( "/pb.OAuthService/CheckToken", "check_token", ctx, request, response); err == nil{
		return response, nil
	}else {
		return nil, err
	}

}


type ClientManager interface {
	DecoratorInvoke(path string, hystrixName string, ctx context.Context, inputVal interface{}, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName string
	loadBalance loadbalance.LoadBalance
	after       []InvokerAfterFunc
	before      []InvokerBeforeFunc
}

type InvokerAfterFunc func() (err error)

type InvokerBeforeFunc func() (err error)

func (manager *DefaultClientManager) DecoratorInvoke(path string, hystrixName string, ctx context.Context, inputVal interface{}, outVal interface{}) (err error) {

			for _, fn := range manager.before{
				if err = fn(); err != nil{
					return err
				}
			}

			if err =  hystrix.Do(hystrixName, func() error {

				instances := discoveryClient.DiscoverServices(manager.serviceName, logger)
				if instance, err := manager.loadBalance.SelectService(instances); err == nil {
					if instance.GrpcPort > 0 {
						if conn, err := grpc.Dial(instance.Host + ":" + strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(), grpc.WithTimeout(1*time.Second)); err == nil {
							if err = conn.Invoke(ctx, path, inputVal, outVal); err != nil{
								return err
							}
						} else {
							return err
						}
					} else {
						return errors.New("no rpc service in " + instance.Host)
					}
				} else {
					return err
				}
				return nil
			}, func(e error) error {
				logger.Println(e.Error())
				return e
			}); err != nil{
				return err
			}else {
				for _, fn := range manager.after{
					if err = fn(); err != nil{
						return err
					}
				}
				return nil
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
		/**
		TODO:需要出示manager，或者直接使用defaultManager
		*/
		manager: &DefaultClientManager{
			serviceName: serviceName,
			loadBalance: lb,
		},
		serviceName: serviceName,
		loadBalance: lb,
	}, nil

}
