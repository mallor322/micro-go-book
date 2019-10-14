package basic

import (
	"github.com/go-kit/kit/sd/consul"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"log"
	"strconv"
	"sync"
)


type ConsulClient interface {

	/**
	 * 服务注册接口
	 * @param serviceName 服务名
	 * @param instanceId 服务实例Id
	 * @param instancePort 服务实例端口
	 * @param healthCheckUrl 健康检查地址
	 * @param meta 服务实例元数据
	 */
	Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool

	/**
	 * 服务注销接口
	 * @param instanceId 服务实例Id
	 */
	DeRegister(instanceId string, logger *log.Logger) bool

	/**
	 * 发现服务实例接口
	 * @param serviceName 服务名
	 */
	DiscoverServices(serviceName string, logger *log.Logger) []interface{}
}

type ConsulClientInstance struct {
	Host string // Consul Host
	Port int 	// Consul Port
	config *api.Config
	client consul.Client
	mutex sync.Mutex
	instancesMap sync.Map}

func New(consulHost string, consulPort int) *ConsulClientInstance{
	// 通过 Consul Host 和 Consul Port 创建一个 consul.Client
	consulConfig := api.DefaultConfig()
	consulConfig.Address = consulHost + ":" +  strconv.Itoa(consulPort)
	apiClient, err := api.NewClient(consulConfig)
	if err != nil{
		return nil
	}

	client := consul.NewClient(apiClient)

	return &ConsulClientInstance{
		Host:consulHost,
		Port:consulPort,
		config:consulConfig,
		client:client,
	}
}

func (consulClient *ConsulClientInstance)Register(serviceName, instanceId, healthCheckUrl string, instanceHost string, instancePort int, meta map[string]string, logger *log.Logger) bool{


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
		logger.Println("Register Service Error!")
		return false
	}
	logger.Println("Register Service Success!")
	return true
}

func (consulClient *ConsulClientInstance) DeRegister(instanceId string, logger *log.Logger) bool {

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
	logger.Println("Deregister Service Success!")

	return true
}

func (consulClient *ConsulClientInstance) DiscoverServices(serviceName string, logger *log.Logger) []interface{} {

	//  该服务已监控并缓存
	instanceList, ok := consulClient.instancesMap.Load(serviceName); if ok{
		instances := make([]interface{}, len(instanceList.([]*api.AgentService)))
		for i := 0; i < len(instances); i++ {
			instances[i] = instanceList.([]*api.AgentService)[i]
		}
		return instances
	}
	// 申请锁
	consulClient.mutex.Lock()
	// 再次检查是否已缓存
	instanceList, ok = consulClient.instancesMap.Load(serviceName); if ok{
		consulClient.mutex.Unlock()
		instances := make([]interface{}, len(instanceList.([]*api.AgentService)))
		for i := 0; i < len(instances); i++ {
			instances[i] = instanceList.([]*api.AgentService)[i]
		}
		return instances
	} else {
		// 注册监控
		go func() {
			params := make(map[string]interface{})
			params["type"] = "service"
			params["service"] = serviceName
			plan, _ := watch.Parse(params)
			plan.Handler = func(u uint64, i interface{}) {
				if i == nil{
					return
				}
				v, ok := i.([]*api.ServiceEntry)
				if !ok || len(v) == 0 {
					return // 数据异常，忽略
				}
				var healthServices []*api.AgentService
				for _, service := range v{
					// 仅保留状态健康的服务实例
					if service.Checks.AggregatedStatus() == api.HealthPassing{
						healthServices = append(healthServices, service.Service)
					}
				}
				consulClient.instancesMap.Store(serviceName, healthServices)
			}
			defer plan.Stop()
			plan.Run(consulClient.config.Address)
		}()
	}
	defer consulClient.mutex.Unlock()

	// 具体查询根据服务名请求服务实例列表并缓存
	entries, _, err := consulClient.client.Service(serviceName, "", false, nil)
	if err != nil{
		logger.Println("Discover Service Error!")
		return nil
	}
	instances := make([]interface{}, len(entries))
	for i := 0; i < len(instances); i++ {
		instances[i] = entries[i].Service
	}
	return instances



}
