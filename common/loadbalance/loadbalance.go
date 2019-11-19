package loadbalance

import (
	"errors"
	"github.com/hashicorp/consul/api"
	"math/rand"
)

// 负载均衡器
type LoadBalance interface {

	SelectService(service [] *api.AgentService) (*api.AgentService, error)

}

var ErrNoInstances = errors.New("service instances are not existed")

type RandomLoadBalance struct {

}
// 随机负载均衡
func (loadBalance *RandomLoadBalance)SelectService(services []*api.AgentService) (*api.AgentService, error) {

	if services == nil || len(services) == 0{
		return nil, ErrNoInstances
	}

	return services[rand.Intn(len(services))], nil
}
