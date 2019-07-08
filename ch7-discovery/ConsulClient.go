package ch7_discovery

import (
	"log"
	"net"
)

type ConsulClient interface {

	/**
	 * 服务注册接口
	 */
	Register(serviceName, instanceId, healthCheckUrl string, instancePort int, meta map[string]string, logger *log.Logger) bool

	/**
	 * 服务注销接口
	 */
	DeRegister(instanceId string, logger *log.Logger) bool


	/**
	 * 发现服务实例接口
	 */
	DiscoverServices(serviceName string) []string

}

func GetLocalIpAddress() string{
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "127.0.0.1"
	}

	for _, address := range addrs {

		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}

		}
	}
	return "127.0.0.1"
}