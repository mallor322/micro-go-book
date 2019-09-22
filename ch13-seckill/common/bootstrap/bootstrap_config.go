package bootstrap

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func init() {
	viper.AutomaticEnv()
	initBootstrapConfig()
	//读取yaml文件
	//v := viper.New()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("err:%s\n", err)
	}
	if err := subParse("http", &HttpConfig); err != nil {
		log.Fatal("Fail to parse Http config", err)
	}
	if err := subParse("discover", &DiscoverConfig); err != nil {
		log.Fatal("Fail to parse Discover config", err)
	}
	if err := subParse("config", &ConfigServerConfig); err != nil {
		log.Fatal("Fail to parse config server", err)
	}
}
func initBootstrapConfig() {
	//设置读取的配置文件
	viper.SetConfigName("bootstrap")
	//添加读取的配置文件路径
	viper.AddConfigPath("ch13-seckill/user-service")
	//windows环境下为%GOPATH，linux环境下为$GOPATH
	viper.AddConfigPath("$GOPATH/src/")
	//设置配置文件类型
	viper.SetConfigType("yaml")
}

/*func main() {
	fmt.Printf(
		"姓名: %s\n爱好: %s\n性别: %s \n年龄: %d \n",
		HttpConfig.Host,
		HttpConfig.Port,
		HttpConfig.ServiceName,
	)
	//反序列化
	//parseYaml(viper.GetViper())

}*/

//func parseYaml(v *viper.Viper, prefix string, cl interface{}) {
//	var httpConfig HttpConf
//	if err := v.Unmarshal(&cl); err != nil {
//		fmt.Printf("err:%s", err)
//	}
//	fmt.Println("resume config:\n ", resumeConfig)
//}
func subParse(key string, value interface{}) error {
	log.Printf("配置文件的前缀为：%v", key)
	sub := viper.Sub(key)
	sub.AutomaticEnv()
	sub.SetEnvPrefix(key)
	return sub.Unmarshal(value)
}
