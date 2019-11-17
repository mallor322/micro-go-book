package loadbalance

import (
	"errors"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"math/rand"
)

// 负载均衡器
type LoadBalance interface {

	SelectService(service [] *discover.ServiceInstance) (*discover.ServiceInstance, error)

}

type RandomLoadBalance struct {

}
// 随机负载均衡
func (loadBalance *RandomLoadBalance)SelectService(services []*discover.ServiceInstance) (*discover.ServiceInstance, error) {

	if services == nil || len(services) == 0{
		return nil, errors.New("service instances are not exist")
	}

	return services[rand.Intn(len(services))], nil
}

type WeightRoundRobinLoadBalance struct {

}

