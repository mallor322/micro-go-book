package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)


var logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
var defaultLoadBalance LoadBalance = &RandomLoadBalance{}
var discoveryClient discover.DiscoveryClient = discover.New("114.67.98.210", "8500")


type OAuthClient interface{

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
		if instance, err := impl.loadBalance.SelectOne(instances); err == nil{

			if rpcPort, ok := instance.Meta["rpcPort"]; ok{
				if conn, err := grpc.Dial(instance.Address + ":" + rpcPort, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second)); err == nil{

					cl := pb.NewOAuthServiceClient(conn)

					resp, err = cl.CheckToken(ctx, request)
				}else {
					return err
				}
			}else {
				return errors.New("no rpc service in " + instance.Address)
			}


		}else {
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
	SelectOne (instances []*api.AgentService) (*api.AgentService, error)
}

type RandomLoadBalance struct {

}

func (* RandomLoadBalance)SelectOne (instances []*api.AgentService) (*api.AgentService, error) {

	if instances == nil || len(instances) == 0{
		return nil, errors.New("no instance existed")
	}

	return instances[rand.Int()%len(instances)], nil

}

func NewOAuthClient(serviceName string, lb LoadBalance) (OAuthClient, error) {

	if lb == nil{
		lb = defaultLoadBalance
	}

	return &OAuthClientImpl{
		serviceName:serviceName,
		loadBalance:lb,
	}, nil

}


func TestOAuthClient(t *testing.T) {
	var (
		grpcAddr = flag.String("addr", ":9008", "gRPC address")
	)
	flag.Parse()
	//tr := localconfig.ZipkinTracer

	//clientTracer := kitzipkin.GRPCClientTrace(tr, kitzipkin.Name("grpc-transport"))
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}

	oauthClient, _ := NewOAuthClient("oauth", nil)

	resp, err := oauthClient.CheckToken(context.Background(), &pb.CheckTokenRequest{
		Token:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGV0YWlscyI6eyJVc2VySWQiOjEsIlVzZXJuYW1lIjoidXNlcm5hbWUiLCJQYXNzd29yZCI6IiIsIkF1dGhvcml0aWVzIjpbIkFkbWluIiwiU3VwZXIiXX0sIkNsaWVudERldGFpbHMiOnsiQ2xpZW50SWQiOiJjbGllbnRJZCIsIkNsaWVudFNlY3JldCI6IiIsIkFjY2Vzc1Rva2VuVmFsaWRpdHlTZWNvbmRzIjoxODAwLCJSZWZyZXNoVG9rZW5WYWxpZGl0eVNlY29uZHMiOjE4MDAwLCJSZWdpc3RlcmVkUmVkaXJlY3RVcmkiOiJodHRwOi8vMTI3LjAuMC4xIiwiQXV0aG9yaXplZEdyYW50VHlwZXMiOlsicGFzc3dvcmQiLCJyZWZyZXNoX3Rva2VuIl19LCJSZWZyZXNoVG9rZW4iOnsiUmVmcmVzaFRva2VuIjpudWxsLCJUb2tlblR5cGUiOiJqd3QiLCJUb2tlblZhbHVlIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlZjMlZ5UkdWMFlXbHNjeUk2ZXlKVmMyVnlTV1FpT2pFc0lsVnpaWEp1WVcxbElqb2lkWE5sY201aGJXVWlMQ0pRWVhOemQyOXlaQ0k2SWlJc0lrRjFkR2h2Y21sMGFXVnpJanBiSWtGa2JXbHVJaXdpVTNWd1pYSWlYWDBzSWtOc2FXVnVkRVJsZEdGcGJITWlPbnNpUTJ4cFpXNTBTV1FpT2lKamJHbGxiblJKWkNJc0lrTnNhV1Z1ZEZObFkzSmxkQ0k2SWlJc0lrRmpZMlZ6YzFSdmEyVnVWbUZzYVdScGRIbFRaV052Ym1Seklqb3hPREF3TENKU1pXWnlaWE5vVkc5clpXNVdZV3hwWkdsMGVWTmxZMjl1WkhNaU9qRTRNREF3TENKU1pXZHBjM1JsY21Wa1VtVmthWEpsWTNSVmNta2lPaUpvZEhSd09pOHZNVEkzTGpBdU1DNHhJaXdpUVhWMGFHOXlhWHBsWkVkeVlXNTBWSGx3WlhNaU9sc2ljR0Z6YzNkdmNtUWlMQ0p5WldaeVpYTm9YM1J2YTJWdUlsMTlMQ0pTWldaeVpYTm9WRzlyWlc0aU9uc2lVbVZtY21WemFGUnZhMlZ1SWpwdWRXeHNMQ0pVYjJ0bGJsUjVjR1VpT2lJaUxDSlViMnRsYmxaaGJIVmxJam9pSWl3aVJYaHdhWEpsYzFScGJXVWlPbTUxYkd4OUxDSmxlSEFpT2pFMU56TTFOemN6TURRc0ltbHpjeUk2SWxONWMzUmxiU0o5LmJneEdtaGJyVmVVQWljMk1ZRWtmeFJFVFVvVHhPYUZmQXFTa2xoQk50N2ciLCJFeHBpcmVzVGltZSI6IjIwMTktMTEtMTNUMDA6NDg6MjQuMzU2OTY4KzA4OjAwIn0sImV4cCI6MTU3MzU2MTEwNCwiaXNzIjoiU3lzdGVtIn0.CRXA_vadztUpOuKPth7HSb-E4l0otZEr6YNU5ZQcEcc",
	})


	defer conn.Close()

	//svr := client.CheckToken(conn, clientTracer)
	//result, err := svr.GetOAuth2DetailsByAccessToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGV0YWlscyI6eyJVc2VySWQiOjEsIlVzZXJuYW1lIjoidXNlcm5hbWUiLCJQYXNzd29yZCI6IiIsIkF1dGhvcml0aWVzIjpbIkFkbWluIiwiU3VwZXIiXX0sIkNsaWVudERldGFpbHMiOnsiQ2xpZW50SWQiOiJjbGllbnRJZCIsIkNsaWVudFNlY3JldCI6IiIsIkFjY2Vzc1Rva2VuVmFsaWRpdHlTZWNvbmRzIjoxODAwLCJSZWZyZXNoVG9rZW5WYWxpZGl0eVNlY29uZHMiOjE4MDAwLCJSZWdpc3RlcmVkUmVkaXJlY3RVcmkiOiJodHRwOi8vMTI3LjAuMC4xIiwiQXV0aG9yaXplZEdyYW50VHlwZXMiOlsicGFzc3dvcmQiLCJyZWZyZXNoX3Rva2VuIl19LCJSZWZyZXNoVG9rZW4iOnsiUmVmcmVzaFRva2VuIjpudWxsLCJUb2tlblR5cGUiOiJqd3QiLCJUb2tlblZhbHVlIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlZjMlZ5UkdWMFlXbHNjeUk2ZXlKVmMyVnlTV1FpT2pFc0lsVnpaWEp1WVcxbElqb2lkWE5sY201aGJXVWlMQ0pRWVhOemQyOXlaQ0k2SWlJc0lrRjFkR2h2Y21sMGFXVnpJanBiSWtGa2JXbHVJaXdpVTNWd1pYSWlYWDBzSWtOc2FXVnVkRVJsZEdGcGJITWlPbnNpUTJ4cFpXNTBTV1FpT2lKamJHbGxiblJKWkNJc0lrTnNhV1Z1ZEZObFkzSmxkQ0k2SWlJc0lrRmpZMlZ6YzFSdmEyVnVWbUZzYVdScGRIbFRaV052Ym1Seklqb3hPREF3TENKU1pXWnlaWE5vVkc5clpXNVdZV3hwWkdsMGVWTmxZMjl1WkhNaU9qRTRNREF3TENKU1pXZHBjM1JsY21Wa1VtVmthWEpsWTNSVmNta2lPaUpvZEhSd09pOHZNVEkzTGpBdU1DNHhJaXdpUVhWMGFHOXlhWHBsWkVkeVlXNTBWSGx3WlhNaU9sc2ljR0Z6YzNkdmNtUWlMQ0p5WldaeVpYTm9YM1J2YTJWdUlsMTlMQ0pTWldaeVpYTm9WRzlyWlc0aU9uc2lVbVZtY21WemFGUnZhMlZ1SWpwdWRXeHNMQ0pVYjJ0bGJsUjVjR1VpT2lJaUxDSlViMnRsYmxaaGJIVmxJam9pSWl3aVJYaHdhWEpsYzFScGJXVWlPbTUxYkd4OUxDSmxlSEFpT2pFMU56TXhOREF5TWpnc0ltbHpjeUk2SWxONWMzUmxiU0o5LjE3bGZrZ3RraFBRVTVkYXA5MGtoQjVKUFRFLXU3V0x4aVZrd2FDcG5uLWsiLCJFeHBpcmVzVGltZSI6IjIwMTktMTEtMDdUMjM6MjM6NDguMTEwNDM1KzA4OjAwIn0sImV4cCI6MTU3MzEyNDAyOCwiaXNzIjoiU3lzdGVtIn0.sOgngpp781LNU6JpCxRCcnYTpZ7YfnAr4-aig29JASo")
	if err != nil {
		fmt.Println("Check error", err.Error())
	}else {
		fmt.Println("result=", resp.ClientDetails.ClientId)
	}
}
