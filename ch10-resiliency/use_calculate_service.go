package main

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/hashicorp/consul/api"
	"github.com/keets2012/Micro-Go-Pracrise/basic"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

func startUseCalculateHttpListener(host string, port int) {
	basic.Server = &http.Server{
		Addr: host + ":" + strconv.Itoa(port),
	}
	// 启动 hystrixStreamHandler 推送统计数据
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	http.Handle("/hystrix/stream", hystrixStreamHandler)
	http.HandleFunc("/health", basic.CheckHealth)
	http.HandleFunc("/use/calculate", useCalculate)
	http.HandleFunc("/discovery", basic.DiscoveryService)
	err := basic.Server.ListenAndServe()
	if err != nil {
		basic.Logger.Println("Service is going to close...")
	}
}

func main() {

	hystrix.ConfigureCommand("Calculate.calculate", hystrix.CommandConfig{
		RequestVolumeThreshold: 4,
	})

	basic.StartService("UseCalculate", "", 10086, startUseCalculateHttpListener)

}

func useCalculate(writer http.ResponseWriter, reader *http.Request) {
	a, _ := strconv.Atoi(reader.URL.Query().Get("a"))
	b, _ := strconv.Atoi(reader.URL.Query().Get("b"))

	result, err := getCalculateResult(a, b)

	if err != nil {
		_, err = fmt.Fprintln(writer, "Get result from Calculate is "+err.Error())
	} else {
		_, err = fmt.Fprintln(writer, "Get result from Calculate is "+result)
	}

	if err != nil {
		basic.Logger.Println(err)
	}
}

func getCalculateResult(a, b int) (string, error) {

	serviceName := "Calculate"

	var result string

	err := hystrix.Do("Calculate.calculate", func() error {

		instances := basic.ConsulService.DiscoverServices(serviceName, basic.Logger)

		if instances == nil || len(instances) == 0 {
			basic.Logger.Println("No Calculate instances are working!")
			return errors.New("No Calculate instances are working")
		}

		// 随机选取一个服务实例进行计算
		rand.Seed(time.Now().UnixNano())
		selectInstance := instances[rand.Intn(len(instances))].(*api.AgentService)

		requestUrl := url.URL{
			Scheme:   "http",
			Host:     selectInstance.Address + ":" + strconv.Itoa(selectInstance.Port),
			Path:     "/calculate",
			RawQuery: "a=" + strconv.Itoa(a) + "&b=" + strconv.Itoa(b),
		}

		resp, err := http.Get(requestUrl.String())
		if err != nil {
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
	} else {
		return "", err
	}
}