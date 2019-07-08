package diy

import (
	"bytes"
	"ch7-discovery"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)


type InstanceInfo struct {

	ID string `json:"ID"`
	Name string `json:"Name"`
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


type ConsulClient struct {
	Host string
	Port int
}


func (consulClient *ConsulClient) Register(serviceName, instanceId, healthCheckUrl string, instancePort int, meta map[string]string, logger *log.Logger) bool{

	instanceHost := ch7_discovery.GetLocalIpAddress()

	instanceInfo := &InstanceInfo{
		ID:      instanceId,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta: meta,
		EnableTagOverride: false,
		Check: Check{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:						"15s",
		},
		Weights: Weights{
			Passing: 10,
			Warning: 1,
		},
	}

	byteData,_ := json.Marshal(instanceInfo)

	req, err := http.NewRequest("PUT",
		"http://" + consulClient.Host + ":" + strconv.Itoa(consulClient.Port) + "/v1/agent/service/register",
		bytes.NewReader(byteData))

	if err == nil {
		req.Header.Set("Content-Type", "application/json;charset=UTF-8")

		client := http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Println("Register Service Error!")
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if body == nil || len(body) == 0 {
				log.Println("Register Service Success!")
				return true;
			} else {
				log.Println("Register Service Error!")
			}
		}
	}
	return false
}

func (consulClient *ConsulClient) DeRegister(instanceId string, logger *log.Logger) bool {

	req, err := http.NewRequest("PUT",
		"http://" + consulClient.Host + ":" + strconv.Itoa(consulClient.Port) + "/v1/agent/service/deregister/" +instanceId, nil)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Deregister Service Error!")
	}else {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if body == nil || len(body) == 0 {
			log.Println("Unregister Service Success!")
			return true
		}else {
			log.Println("Deregister Service Error!")
		}
	}
	return false
}

func New(consulHost string, consulPort int) *ConsulClient {
	return &ConsulClient{
		Host: consulHost,
		Port: consulPort,
	}
}

func (consulClient *ConsulClient) DiscoverServices(serviceName string) []string{


	req, err := http.NewRequest("GET",
		"http://" + consulClient.Host + ":" + strconv.Itoa(consulClient.Port) + "/v1/agent/health/service/name/" + serviceName, nil)

	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Println("Discover Service Error!")
	}else {

		var serviceList [] struct {
			Service InstanceInfo `json:"Service"`
		}
		err = json.NewDecoder(resp.Body).Decode(&serviceList)
		resp.Body.Close()
		if err == nil {
			instances := make([]string, len(serviceList))
			for i := 0; i < len(instances); i++ {
				instance, err := json.Marshal(serviceList[i].Service)
				if err == nil {
					instances[i] = string(instance)
				}
			}
			return instances
		}
	}
	return nil
}






