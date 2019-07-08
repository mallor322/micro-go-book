package main

import (
	"ch7-discovery"
	"ch7-discovery/diy"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
)

var server *http.Server


func sayHello(writer http.ResponseWriter, reader *http.Request)  {
	_, err := fmt.Fprintln(writer, "Hello World!")
	if err != nil{
		logger.Println(err)
	}
}

func startHttpListener(port int)  {
	server = &http.Server{
		Addr: ch7_discovery.GetLocalIpAddress() + ":" +strconv.Itoa(port),
	}
	http.HandleFunc("/health", CheckHealth)
	http.HandleFunc("/sayHello", sayHello)
	http.HandleFunc("/discovery", discoveryService)

	err := server.ListenAndServe()

	if err != nil{
		logger.Println("Service is going to close...")
	}

}

func CheckHealth(writer http.ResponseWriter, reader *http.Request)  {
	logger.Println("Health check starts!")
	_, err := fmt.Fprintln(writer, "Server is OK!")
	if err != nil{
		logger.Println(err)
	}
}

func discoveryService(writer http.ResponseWriter, reader *http.Request)  {
	serviceName := reader.URL.Query().Get("serviceName")
	instances := consulClient.DiscoverServices(serviceName)
	writer.Header().Set("Content-Type", "application/json")
	jsonRes, _ := json.Marshal(instances)
	_, err := writer.Write(jsonRes)
	if err != nil{
		logger.Println(err)
	}

}


func closeServer( waitGroup *sync.WaitGroup, exit <-chan os.Signal, instanceId string, logger *log.Logger)  {


	<- exit
	waitGroup.Add(1)
	consulClient.DeRegister(instanceId, logger)

	err := server.Shutdown(nil)
	if err != nil{
		log.Println(err)
	}
	waitGroup.Done()

}


var consulClient ch7_discovery.ConsulClient
var logger *log.Logger

func main()  {


	consulClient = diy.New("10.224.19.186", 8500)

	if consulClient == nil{
		panic(0)
	}

	instanceId := uuid.NewV4().String()

	logger = log.New(os.Stderr, "", log.LstdFlags)

	if !consulClient.Register("SayHello", instanceId, "/health", 10086, nil, logger) {
		// 注册失败，服务启动失败
		panic(0)
	}

	exit := make(chan os.Signal)
	// 仅监控 ctrl + c
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	var waitGroup sync.WaitGroup

	go closeServer(&waitGroup, exit, instanceId, logger)

	// 在主线程启动http服务器
	startHttpListener(10086)

	waitGroup.Wait()

	log.Println("Closed the Server!")

}
