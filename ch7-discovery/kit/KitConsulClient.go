package kit

import (
	ch7discovery "ch7-discovery"
	"encoding/json"
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"log"
	"strconv"
)

type ConsulClient struct {
	Host string
	Port int
	client consul.Client
}


func New(consulHost string, consulPort int) *ConsulClient{

	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" +  strconv.Itoa(consulPort)

	apiClient, err := api.NewClient(consulConfig)
	if err != nil{
		return nil
	}

	client := consul.NewClient(apiClient)

	return &ConsulClient{
		Host:consulHost,
		Port:consulPort,
		client:client,
	}
}

func (consulClient *ConsulClient)Register(serviceName, instanceId, healthCheckUrl string, instancePort int, meta map[string]string, logger *log.Logger) bool{

	instanceHost := ch7discovery.GetLocalIpAddress()

	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    serviceName,
		Address: instanceHost,
		Port:    instancePort,
		Meta: meta,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + instanceHost + ":" + strconv.Itoa(instancePort) + healthCheckUrl,
			Interval:						"15s",
		},
	}
	err := consulClient.client.Register(serviceRegistration)

	if err != nil{
		log.Println("Register Service Error!")
		return false
	}
	log.Println("Register Service Success!")

	return true
}

func (consulClient *ConsulClient) DeRegister(instanceId string, logger *log.Logger) bool {

	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
	}
	err := consulClient.client.Deregister(serviceRegistration)

	if err != nil{
		logger.Println("Deregister Service Error!")
		return false
	}
	log.Println("Deregister Service Success!")

	return true
}

func (consulClient *ConsulClient) DiscoverServices(serviceName string) []string {

	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil{
		log.Println("Discover Service Error!")
		return nil
	}

	instances := make([]string, len(entries))
	for i := 0; i < len(instances); i++ {
		instance, err := json.Marshal(entries[i].Service)
		if err == nil {
			instances[i] = string(instance)
		}
	}
	return instances
}