package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/loadbalance"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
)

var logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
var defaultLoadBalance loadbalance.LoadBalance = &loadbalance.RandomLoadBalance{}
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
	loadBalance loadbalance.LoadBalance
}

func (impl *OAuthClientImpl) CheckToken(request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {
	var funcAfterDecor1 = impl.CheckTokenInternale
	impl.manager.Decorator(&funcAfterDecor1, OAuthClient.CheckToken, "/pb.OAuthService/CheckToken", context.Background(), &pb.CheckTokenResponse{})
	return funcAfterDecor1(&pb.CheckTokenRequest{})
}

func (impl *OAuthClientImpl) CheckTokenInternale(request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {

	return nil, nil

}

type ClientManager interface {
	Decorator(decoPtr, fn interface{}, path string, ctx context.Context, outVal interface{}) (err error)
}

type DefaultClientManager struct {
	serviceName string
	loadBalance loadbalance.LoadBalance
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

			hystrix.Do("test_check_token", func() error {

				instances := discoveryClient.DiscoverServices("serviceName", logger)
				if instance, err := manager.loadBalance.SelectService(instances); err == nil {

					if instance.GrpcPort > 0 {
						if conn, err := grpc.Dial(instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(), grpc.WithTimeout(1*time.Second)); err == nil {
							// TODO: in需要扩充
							// 如何使用反射？？？new出这个instance？
							conn.Invoke(ctx, path, in[0], outVal, nil)
							out[0] = reflect.ValueOf(outVal)
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
			})

			fmt.Println("before")

			out = targetFunc.Call(in)
			fmt.Println("after")
			return
		})

	decoratedFunc.Set(v)
	return
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
