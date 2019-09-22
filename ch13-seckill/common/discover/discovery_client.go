package discover

import (
	"github.com/go-kit/kit/sd/consul"
	"log"
)

type DiscoveryClientInstance struct {
	Host   string //  Host
	Port   int    //  Port
	client consul.Client
}

type ServiceInstance struct {
	Host     string //  Host
	Port     int    //  Port
	GrpcPort int
}
type DiscoveryClient interface {
	/**
	 * 服务注册接口
	 * @param serviceName 服务名
	 * @param instanceId 服务实例Id
	 * @param instancePort 服务实例端口
	 * @param healthCheckUrl 健康检查地址
	 * @param meta 服务实例元数据
	 */
	Register(instanceId, svcHost, healthCheckUrl, svcPort string, svcName string, meta map[string]string, tags []string, logger *log.Logger) bool

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
