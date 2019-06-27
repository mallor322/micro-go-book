package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type InstanceInfo struct {

	ID string `json:"ID"`
	Name string `json:"name"`
	Tags []string `json:"Tags,omitempty"`
	Address string `json:"Address"`
	Port int `json:"Port"`
	Meta map[string]string `json:"Meta,omitempty"`
	EnableTagOverride bool `json:"EnableTagOverride"`
	Check `json:"Check,omitempty"`
	Weights `json:"Weights,omitempty"`

}

type Check struct {

	DeregisterCriticalServiceAfter string `json:"DeregisterCriticalServiceAfter"`
	Args []string `json:"Args,omitempty"`
	HTTP string `json:"HTTP"`
	Interval string `json:"Interval,omitempty"`
	TTL string `json:"TTL,omitempty"`

}

type Weights struct {

	Passing int `json:"Passing"`
	Warning int `json:"Warning"`

}




func CheckHealth(writer http.ResponseWriter, reader *http.Request)  {
	fmt.Print("Health check starts at ")
	fmt.Println(time.Now())
	fmt.Fprintln(writer, "Server is OK!")
}


var consulAddress = "http://10.224.19.186:8500/";
var registerUrl = consulAddress + "v1/agent/service/register"
var unregisterUrl = consulAddress + "v1/agent/service/deregister/"
var healthServiceUrl = consulAddress + "v1/agent/health/service/name/";



func RegisterService(instanceInfo *InstanceInfo) bool{

	if instanceInfo == nil {
		fmt.Println("InstanceInfo could not be NIL!")
		return false
	}

	byteData,_ := json.Marshal(instanceInfo)

	req, err := http.NewRequest("PUT",
		registerUrl,
		bytes.NewReader(byteData))

	req.Header.Set("Content-Type", "application/json;charset=UTF-8")

	client := http.Client{}

	resp, err := client.Do(req)


	if err != nil {
		fmt.Println("Register Service Error!")
	}else {
		body, _ := ioutil.ReadAll(resp.Body)
		if body == nil || len(body) == 0 {
			fmt.Println("Register Service Success!")
			return true;
		}else {
			fmt.Println("Register Service Error!")
		}
	}
	return false
}

func UnregisterService(serviceId string) bool {

	req, err := http.NewRequest("PUT",
		unregisterUrl + serviceId, nil)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Unregister Service Error!")
	}else {
		body, _ := ioutil.ReadAll(resp.Body)
		if body == nil || len(body) == 0 {
			fmt.Println("Unregister Service Success!")
			return true
		}else {
			fmt.Println("unregister Service Error!")
		}
	}
	return false
}


func QueryInstanceListByServiceName(serviceName string) []InstanceInfo{


	req, err := http.NewRequest("GET",
		healthServiceUrl + serviceName, nil)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Unregister Service Error!")
	}else {
		body, _ := ioutil.ReadAll(resp.Body)
		if body == nil || len(body) == 0 {
			return nil
		}else {
			var serviceList [] struct {
				Service InstanceInfo `json:"Service"`
			}
			json.Unmarshal(body, &serviceList)
			serviceInfos := make([]InstanceInfo, len(serviceList))
			for i:= 0 ; i < len(serviceInfos) ; i++{
				serviceInfos[i] = serviceList[i].Service;
			}
			return serviceInfos
		}
	}
	return nil
}







