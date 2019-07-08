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
	Host string // Consul Host
	Port int 	// Consul Port
	client consul.Client
}

func New(consulHost string, consulPort int) *ConsulClient{
	// 通过 Consul Host 和 Consul Port 创建一个 consul.Client
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

	// 获取服务实例 IP
	instanceHost := ch7discovery.GetLocalIpAddress()

	// 1. 构建服务实例元数据
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

	// 2. 发送服务注册到 Consul 中
	err := consulClient.client.Register(serviceRegistration)

	if err != nil{
		log.Println("Register Service Error!")
		return false
	}
	log.Println("Register Service Success!")
	return true
}

func (consulClient *ConsulClient) DeRegister(instanceId string, logger *log.Logger) bool {

	// 构建包含服务实例 ID 的元数据结构体
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
	}
	// 发送服务注销请求
	err := consulClient.client.Deregister(serviceRegistration)

	if err != nil{
		logger.Println("Deregister Service Error!")
		return false
	}
	log.Println("Deregister Service Success!")

	return true
}

func (consulClient *ConsulClient) DiscoverServices(serviceName string) []interface{} {

	// 根据服务名请求服务实例列表，可以添加额外的筛选参数
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil{
		log.Println("Discover Service Error!")
		return nil
	}

	instances := make([]interface{}, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = entries[i].Service
	}
	return instances
}