package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var server *http.Server


func sayHello(writer http.ResponseWriter, reader *http.Request)  {
	fmt.Fprintln(writer, "Hello World!")
}

func startHttpListener()  {

	server = &http.Server{
		Addr: "10.93.246.254:10088",
	}

	http.HandleFunc("/health", CheckHealth)
	http.HandleFunc("/sayHello", sayHello)
	http.HandleFunc("/discovery", discoveryService)

	server.ListenAndServe()

}

func discoveryService(writer http.ResponseWriter, reader *http.Request)  {
	serviceName := reader.URL.Query().Get("serviceName")
	instances := QueryInstanceListByServiceName(serviceName)
	writer.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(&instances)
	writer.Write(json)

}


func closeServer( waitGroup *sync.WaitGroup, exit <-chan os.Signal)  {

	<- exit
	waitGroup.Add(1)
	UnregisterService("OKBUG001")

	err := server.Shutdown(nil)
	if err != nil{
		fmt.Println(err)
	}

	waitGroup.Done()

}


func main()  {


	if !RegisterService(&InstanceInfo{
		ID:      "OKBUG001",
		Name:    "OKBUG",
		Address: "10.93.246.254",
		Port:    10086,
		Meta: map[string]string{
			"feature": "GOLD",
			"version": "Newest",
		},
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://10.93.246.254:10088/health",
			Interval:						"15s",
		},
		Weights: Weights{
			Passing: 10,
			Warning: 1,
		},
	}) {
		// 注册失败，服务启动失败
		panic(0)
	}

	exit := make(chan os.Signal)
	// 仅监控 ctrl + c
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	var waitGroup sync.WaitGroup

	go closeServer(&waitGroup, exit)

	// 在主线程启动http服务器
	startHttpListener()

	waitGroup.Wait()

	fmt.Println("Closed the Server!")


}
