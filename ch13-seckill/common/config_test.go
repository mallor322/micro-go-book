package main

import (
	"fmt"
	conf "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/config"
	"github.com/spf13/viper"
	"testing"
)

func TestConfig(t *testing.T) {
	fmt.Printf("姓名：%s，\n性别：%s，\n年龄 %s!", viper.GetString("discover.port"), conf.HttpConfig.ServiceName, conf.DiscoverConfig.Host) //这个写入到w的是输出到客户端的
}
