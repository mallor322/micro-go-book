package main

import (
	"fmt"
	register "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"testing"
)

func TestDiscover(t *testing.T) {
	register.Register()
	ins := register.DiscoveryService("config-service")
	fmt.Printf("selected instance's host %s, port %d\n", ins.Host, ins.Port)
	register.Deregister()
}
