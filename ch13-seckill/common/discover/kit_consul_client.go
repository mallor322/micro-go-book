package discover

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"log"
	"strconv"
)

func New(consulHost string, consulPort string) *DiscoveryClientInstance {
	port, _ := strconv.Atoi(consulPort)
	// 通过 Consul Host 和 Consul Port 创建一个 consul.Client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" + strconv.Itoa(port)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil {
		return nil
	}

	client := consul.NewClient(apiClient)

	return &DiscoveryClientInstance{
		Host:   consulHost,
		Port:   port,
		client: client,
	}
}

func (consulClient *DiscoveryClientInstance) Register(instanceId, svcHost, healthCheckUrl, svcPort string, svcName string, meta map[string]string, tags []string, logger *log.Logger) bool {
	port, _ := strconv.Atoi(svcPort)

	// 1. 构建服务实例元数据
	serviceRegistration := &api.AgentServiceRegistration{
		ID:      instanceId,
		Name:    svcName,
		Address: svcHost,
		Port:    port,
		Meta:    meta,
		Tags:    tags,
		Check: &api.AgentServiceCheck{
			DeregisterCriticalServiceAfter: "30s",
			HTTP:                           "http://" + svcHost + ":" + strconv.Itoa(port) + healthCheckUrl,
			Interval:                       "15s",
		},
	}

	// 2. 发送服务注册到 Consul 中
	err := consulClient.client.Register(serviceRegistration)

	if err != nil {
		logger.Println("Register Service Error!")
		return false
	}
	logger.Println("Register Service Success!")
	return true
}

func (consulClient *DiscoveryClientInstance) DeRegister(instanceId string, logger *log.Logger) bool {

	// 构建包含服务实例 ID 的元数据结构体
	serviceRegistration := &api.AgentServiceRegistration{
		ID: instanceId,
	}
	// 发送服务注销请求
	err := consulClient.client.Deregister(serviceRegistration)

	if err != nil {
		logger.Println("Deregister Service Error!")
		return false
	}
	logger.Println("Deregister Service Success!")

	return true
}

func (consulClient *DiscoveryClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {

	// 根据服务名请求服务实例列表，可以添加额外的筛选参数
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil {
		logger.Println("Discover Service Error!")
		return nil
	}

	instances := make([]interface{}, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = entries[i].Service
	}
	return instances
}
