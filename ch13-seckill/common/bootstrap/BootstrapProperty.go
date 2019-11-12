package bootstrap

var (
	HttpConfig         HttpConf
	DiscoverConfig     DiscoverConf
	ConfigServerConfig ConfigServerConf
	RpcConfig		   RpcConf
)

//Http配置
type HttpConf struct {
	Host        string
	Port        string
	ServiceName string
}

// RPC配置
type RpcConf struct {
	Port	string
	Host 	string
}

//服务发现与注册配置
type DiscoverConf struct {
	Host       string
	Port       string
	InstanceId string
}

//配置中心
type ConfigServerConf struct {
	Id      string
	Profile string
	Label   string
}
