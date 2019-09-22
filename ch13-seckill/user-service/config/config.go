package config

import (
	_ "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	conf "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/config"
	"github.com/spf13/viper"
	"log"
)

const (
	kConfigType = "CONFIG_TYPE"
)

func init() {
	viper.AutomaticEnv()
	initDefault()

	if err := conf.LoadRemoteConfig(); err != nil {
		log.Fatal("Fail to load remote config", err)
	}

	if err := conf.Sub("mysql", &conf.MysqlConfig); err != nil {
		log.Fatal("Fail to parse config", err)
	}
}

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")
}
