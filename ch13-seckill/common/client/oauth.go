package client

import (
	"context"
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"os"
	"time"
)

var logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
var defaultLoadBalance LoadBalance = &RandomLoadBalance{}
var discoveryClient discover.DiscoveryClient = discover.New("114.67.98.210", "8500")

type OAuthClient interface {
	CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error)
}

type OAuthClientImpl struct {
	serviceName string
	loadBalance LoadBalance
}

func (impl *OAuthClientImpl) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {

	var resp *pb.CheckTokenResponse

	err := hystrix.Do("test_check_token", func() error {

		instances := discoveryClient.DiscoverServices(impl.serviceName, logger)
		if instance, err := impl.loadBalance.SelectOne(instances); err == nil {

			if rpcPort, ok := instance.Meta["rpcPort"]; ok {
				if conn, err := grpc.Dial(instance.Address+":"+rpcPort, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second)); err == nil {

					cl := pb.NewOAuthServiceClient(conn)

					resp, err = cl.CheckToken(ctx, request)
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

	return resp, err

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
		serviceName: serviceName,
		loadBalance: lb,
	}, nil

}
