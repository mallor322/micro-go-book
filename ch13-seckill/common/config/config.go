package conf

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/discover"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strconv"
)

const (
	kConfigType = "CONFIG_TYPE"
)

func init2() {
	viper.AutomaticEnv()
	initDefault()

	if err := LoadRemoteConfig(); err != nil {
		log.Fatal("Fail to load config", err)
	}

	if err := Sub("redis", &Redis); err != nil {
		log.Fatal("Fail to parse config", err)
	}
	if err := Sub("etcd", &Etcd); err != nil {
		log.Fatal("Fail to parse config", err)
	}
	if err := Sub("service", &SecKill); err != nil {
		log.Fatal("Fail to parse config", err)
	}
}

func initDefault() {
	viper.SetDefault(kConfigType, "yaml")

}

func LoadRemoteConfig() (err error) {
	serviceInstance := discover.DiscoveryService(bootstrap.ConfigServerConfig.Id)
	configServer := "http://" + serviceInstance.Host + ":" + strconv.Itoa(serviceInstance.Port)
	confAddr := fmt.Sprintf("%v/%v/%v-%v.%v",
		configServer, bootstrap.ConfigServerConfig.Label,
		bootstrap.HttpConfig.ServiceName, bootstrap.ConfigServerConfig.Profile,
		viper.Get(kConfigType))
	resp, err := http.Get(confAddr)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	viper.SetConfigType(viper.GetString(kConfigType))
	if err = viper.ReadConfig(resp.Body); err != nil {
		return
	}
	log.Println("Load config from: ", confAddr)
	return
}

func Sub(key string, value interface{}) error {
	log.Printf("配置文件的前缀为：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
