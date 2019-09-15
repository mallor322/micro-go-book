package basic

import (
	"encoding/json"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var Server *http.Server



//
func CheckHealth(writer http.ResponseWriter, reader *http.Request)  {
	Logger.Println("Health check!")
	_, err := fmt.Fprintln(writer, "Server is OK!")
	if err != nil{
		Logger.Println(err)
	}
}

func DiscoveryService(writer http.ResponseWriter, reader *http.Request)  {
	serviceName := reader.URL.Query().Get("serviceName")
	instances := ConsulService.DiscoverServices(serviceName, Logger)
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(writer).Encode(instances)
	if err != nil{
		Logger.Println(err)
	}
}


func closeServer( waitGroup *sync.WaitGroup, exit <-chan os.Signal, instanceId string, logger *log.Logger)  {
	// 等待关闭信息通知
	<- exit
	// 主线程等待
	waitGroup.Add(1)
	// 服务注销
	ConsulService.DeRegister(instanceId, logger)
	// 关闭 http 服务器
	err := Server.Shutdown(nil)
	if err != nil{
		log.Println(err)
	}
	// 主线程可继续执行
	waitGroup.Done()
}

func StartService(serviceName string, host string, port int, serviceFunc func(host string, port int))  {

	// 1.实例化一个 Consul 客户端，此处实例化了原生态实现版本
	ConsulService = New("127.0.0.1", 8500)
	//// 实例失败，停止服务
	if ConsulService == nil{
		panic(0)
	}

	// 通过 go.uuid 获取一个服务实例ID
	instanceId := uuid.NewV4().String()
	Logger = log.New(os.Stderr, "", log.LstdFlags)
	// 服务注册
	if !ConsulService.Register(serviceName, instanceId, "/health", host, port, nil, Logger) {
		// 注册失败，服务启动失败
		panic(0)
	}

	// 2.建立一个通道监控系统信号
	exit := make(chan os.Signal)
	// 仅监控 ctrl + c
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	var waitGroup sync.WaitGroup
	// 注册关闭事件，等待 ctrl + c 系统信号通知服务关闭
	go closeServer(&waitGroup, exit, instanceId, Logger)

	// 3. 在主线程启动http服务器
	serviceFunc(host, port)

	// 等待关闭事件执行结束，结束主线程
	waitGroup.Wait()
	log.Println("Closed the Server!")

}

var ConsulService ConsulClient
var Logger *log.Logger

