package hystrix

import (
	"context"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	"github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"io"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	commandName := "my-endpoint"
	hystrix.ConfigureCommand(commandName, hystrix.CommandConfig{
		Timeout:                1000 * 30,
		ErrorPercentThreshold:  1,
		SleepWindow:            10000,
		MaxConcurrentRequests:  1000,
		RequestVolumeThreshold: 5,
	})
	port, _ := strconv.Atoi(consulPort)
	// 通过 Consul Host 和 Consul Port 创建一个 consul.Client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(port)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return
	}

	client := consul.NewClient(apiClient)
	logger := log.NewNopLogger()
	instancer := consul.NewInstancer(client, logger, "user_service", []string{}, false)
	// 创建hystric 的enpoint.Middleware
	breakerMw := circuitbreaker.Hystrix(commandName)
	//创建端点管理器， 此管理器根据Factory和监听的到实例创建endPoint并订阅instancer的变化动态更新Factory创建的endPoint
	endpointer := sd.NewEndpointer(instancer, reqFactory, logger)
	//创建负载均衡器
	balancer := lb.NewRoundRobin(endpointer)
	reqEndPoint := lb.Retry(3, 100*time.Second, balancer)
	reqEndPoint = breakerMw(reqEndPoint)
	//增加熔断中间件
	reqEndPoint = breakerMw(reqEndPoint)
	//现在我们可以通过 endPoint 发起请求了
	req := struct{}{}
	ctx := context.Background()
	for i := 1; i <= 20; i++ {
		if _, err = reqEndPoint(ctx, req); err != nil {
			fmt.Println("当前时间: ", time.Now().Format("2006-01-02 15:04:05.99"))
			fmt.Println(err)
			//time.sleep(1 * time.Second)
		}
	}
}

//通过传入的 实例地址  创建对应的请求endPoint
func reqFactory(instanceAddr string) (endpoint.Endpoint, io.Closer, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fmt.Println("请求服务: ", instanceAddr, "当前时间: ", time.Now().Format("2006-01-02 15:04:05.99"))
		conn, err := grpc.Dial(instanceAddr, grpc.WithInsecure())
		if err != nil {
			fmt.Println(err)
			panic("connect error")
		}
		defer conn.Close()
		bookClient := book.NewBookServiceClient(conn)
		bi, err := bookClient.GetBookInfo(context.Background(), &book.BookInfoParams{BookId: 1})
		fmt.Println(bi)
		fmt.Println("       ", "获取书籍详情")
		fmt.Println("       ", "bookId: 1", " => ", "bookName:", bi.BookName)
		return nil, nil
	}, nil, nil

}
