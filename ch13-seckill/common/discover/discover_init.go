package discover

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
)

var ConsulService DiscoveryClient
var Logger *log.Logger

func init() {
	// 1.实例化一个 Consul 客户端，此处实例化了原生态实现版本
	ConsulService = New(bootstrap.DiscoverConfig.Host, bootstrap.DiscoverConfig.Port)
	Logger = log.New(os.Stderr, "", log.LstdFlags)

}

//
func CheckHealth(writer http.ResponseWriter, reader *http.Request) {
	Logger.Println("Health check!")
	_, err := fmt.Fprintln(writer, "Server is OK!")
	if err != nil {
		Logger.Println(err)
	}
}

func DiscoveryService(serviceName string) ServiceInstance {
	instances := ConsulService.DiscoverServices(serviceName, Logger)
	//todo lb from config center, default is random.
	if len(instances) < 1 {
		Logger.Printf("no available client for %s.", serviceName)
		//todo 异常处理机制
		os.Exit(1)
	}
	selectOne := instances[0]

	return ServiceInstance{
		Host:     selectOne.Address,
		Port:     selectOne.Port,
		GrpcPort: selectOne.Port - 1,
	}
}

func Register() {
	//// 实例失败，停止服务
	if ConsulService == nil {
		panic(0)
	}

	//判空 instanceId,通过 go.uuid 获取一个服务实例ID
	instanceId := bootstrap.DiscoverConfig.InstanceId

	if instanceId == "" {
		instanceId = bootstrap.HttpConfig.ServiceName + uuid.NewV4().String()
	}

	if !ConsulService.Register(instanceId, bootstrap.HttpConfig.Host, "/health",
		bootstrap.HttpConfig.Port, bootstrap.HttpConfig.ServiceName, nil, nil, Logger) {
		Logger.Printf("register for service %s failed.", bootstrap.HttpConfig.ServiceName)
		// 注册失败，服务启动失败
		panic(0)
	}
	Logger.Printf("register for service %s success.", bootstrap.HttpConfig.ServiceName)

}

func Deregister() {
	//// 实例失败，停止服务
	if ConsulService == nil {
		panic(0)
	}
	//判空 instanceId,通过 go.uuid 获取一个服务实例ID
	instanceId := bootstrap.DiscoverConfig.InstanceId

	if instanceId == "" {
		instanceId = bootstrap.HttpConfig.ServiceName + uuid.NewV4().String()
	}
	if !ConsulService.DeRegister(instanceId, Logger) {
		Logger.Printf("deregister for service %s failed.", bootstrap.HttpConfig.ServiceName)
		panic(0)
	}
}
