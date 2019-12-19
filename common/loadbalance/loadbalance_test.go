package loadbalance

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pkg/discover"
	"testing"
)

func TestMain(m *testing.M) {
	instances := make([]*discover.ServiceInstance, 10)
	for i := 0; i < len(instances); i++ {
		instances[i] = &discover.ServiceInstance{
			Host:      string(i),
			Port:      i,
			GrpcPort:  (i + 1),
			Weight:    i,
			CurWeight: 0,
		}
	}

	load := &WeightRoundRobinLoadBalance{}
	for i := 0; i < 2*len(instances); i++ {
		service, _ := load.SelectService(instances)
		fmt.Println("service is : ", service.Host)
	}

}
