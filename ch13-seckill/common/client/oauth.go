package client

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
	serviceName string
	loadBalance loadbalance.LoadBalance
}

func (impl *OAuthClientImpl) CheckToken(ctx context.Context, request *pb.CheckTokenRequest) (*pb.CheckTokenResponse, error) {

	var resp *pb.CheckTokenResponse

	err := hystrix.Do("test_check_token", func() error {

		instances := discoveryClient.DiscoverServices(impl.serviceName, logger)
		if instance, err := impl.loadBalance.SelectService(instances); err == nil {

			if instance.GrpcPort > 0{
				if conn, err := grpc.Dial(instance.Host+":"+strconv.Itoa(instance.GrpcPort), grpc.WithInsecure(), grpc.WithTimeout(1*time.Second)); err == nil {
					cl := pb.NewOAuthServiceClient(conn)
					resp, err = cl.CheckToken(ctx, request)
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

	return resp, err

}




func NewOAuthClient(serviceName string, lb loadbalance.LoadBalance) (OAuthClient, error) {
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
