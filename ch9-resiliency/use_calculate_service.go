package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func startUseCalculateHttpListener(host string, port int)  {
	server = &http.Server{
		// GetLocalIpAddress用于获取本地IP，可以手动写入
		Addr: host + ":" +strconv.Itoa(port),
	}
	http.HandleFunc("/health", checkHealth)
	http.HandleFunc("/use/calculate", useCalculate)
	http.HandleFunc("/discovery", discoveryService)
	err := server.ListenAndServe()
	if err != nil{
		logger.Println("Service is going to close...")
	}
}

func main()  {

	hystrix.ConfigureCommand("Calculate.calculate", hystrix.CommandConfig{
		RequestVolumeThreshold:4,
	})

	startService("UseCalculate", "127.0.0.1", 10086, startUseCalculateHttpListener)

}


func useCalculate(writer http.ResponseWriter, reader *http.Request)  {
	a, _:= strconv.Atoi(reader.URL.Query().Get("a"))
	b, _:= strconv.Atoi(reader.URL.Query().Get("b"))

	result, err := getCalculateResult(a, b)

	if err != nil{
		_, err = fmt.Fprintln(writer, "Get result from Calculate is " + err.Error())
	}else {
		_, err = fmt.Fprintln(writer, "Get result from Calculate is " + result)
	}

	if err != nil{
		logger.Println(err)
	}
}




func getCalculateResult(a, b int) (string, error) {

	serviceName := "Calculate"

	var result string

	err := hystrix.Do("Calculate.calculate", func() error{

		instances := consulClient.DiscoverServices(serviceName)

		if instances == nil || len(instances) == 0 {
			logger.Println("No Calculate instances are working!")
			return errors.New("No Calculate instances are working")
		}

		// 随机选取一个服务实例进行计算
		rand.Seed(time.Now().UnixNano())
		selectInstance := instances[rand.Intn(len(instances))].(*api.AgentService)

		requestUrl := url.URL{
			Scheme:"http",
			Host:selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
			Path:"/calculate",
			RawQuery:"a=" + strconv.Itoa(a) + "&b=" + strconv.Itoa(b),
		}

		resp, err := http.Get(requestUrl.String())
		if err != nil{
			return err
		}
		body, _ := ioutil.ReadAll(resp.Body)
		result = string(body)

		return nil

	}, func(e error) error {
		return errors.New("Http errors！")
	})

	if err == nil {
		return result, nil
	}else {
		return "", err
	}
}


