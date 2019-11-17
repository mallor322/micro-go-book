package other_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"os"
	"reflect"
	"time"
)

var logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
var defaultLoadBalance LoadBalance = &RandomLoadBalance{}
var discoveryClient discover.DiscoveryClient = discover.New("114.67.98.210", "8500")

type OAuthClient interface {
	CheckToken(request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	/**
	* 可以配置负载均衡策略，重试、等机制。也可以配置invokeAfter和invokerBefore
	 */
	manager     ClientManager
	serviceName string
	loadBalance LoadBalance
}

func (impl *OAuthClientImpl) CheckToken(request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	var funcAfterDecor1 = impl.CheckToken
	impl.manager.Decorator(&funcAfterDecor1, OAuthClient.CheckToken, "/pb.OAuthService/CheckToken", context.Background(), &pb.CheckTokenResponse{})
	return funcAfterDecor1(&pb.CheckTokenRequest{})
}

type ClientManager interface {
	Decorator(decoPtr, fn interface{}, path string, ctx context.Context, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName string
	loadBalance LoadBalance
	after       []InvokerAfterFunc
	before      []InvokerBeforeFunc
}

type InvokerAfterFunc func() (err error)

type InvokerBeforeFunc func() (err error)

func (manager *DefaultClientManager) Decorator(decoPtr, fn interface{}, path string, ctx context.Context, outVal interface{}) (err error) {
	var decoratedFunc, targetFunc reflect.Value

	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)
	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {

			err := hystrix.Do("test_check_token", func() error {

				instances := discoveryClient.DiscoverServices("serviceName", logger)
				if instance, err := manager.loadBalance.SelectOne(instances); err == nil {

					if rpcPort, ok := instance.Meta["rpcPort"]; ok {
						if conn, err := grpc.Dial(instance.Address+":"+rpcPort, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second)); err == nil {
							// TODO: in需要扩充
							// 如何使用反射？？？new出这个instance？
							_ := conn.Invoke(ctx, path, in[0], outVal, nil)
							out[0] = reflect.ValueOf(outVal)
						} else {
							return err
						}
					} else {
						return errors.New("no rpc service in " + instance.Address)
					}

				} else {
					return err
				}
				return nil
			}, func(e error) error {
				logger.Println(e.Error())
				return e
			})

			fmt.Println("before")

			out = targetFunc.Call(in)
			fmt.Println("after")
			return
		})

	decoratedFunc.Set(v)
	return
}

type LoadBalance interface {
	SelectOne(instances []*api.AgentService) (*api.AgentService, error)
}

type RandomLoadBalance struct {
}

func (*RandomLoadBalance) SelectOne(instances []*api.AgentService) (*api.AgentService, error) {

	if instances == nil || len(instances) == 0 {
		return nil, errors.New("no instance existed")
	}

	return instances[rand.Int()%len(instances)], nil

}

func NewOAuthClient(serviceName string, lb LoadBalance) (OAuthClient, error) {
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
		serviceName: serviceName,
		loadBalance: lb,
	}, nil

}
