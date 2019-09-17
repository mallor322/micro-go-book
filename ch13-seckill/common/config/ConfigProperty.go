package conf

import (
	"github.com/go-redis/redis"
	"go.etcd.io/etcd/clientv3"
	"sync"
)

var (
	Redis RedisConf
	Etcd EtcdConf
	SecKill SecKillConf
)

type EtcdConf struct {
	EtcdConn          *clientv3.Client //链接
	EtcdSecProductKey string           //商品键
}

//redis配置
type RedisConf struct {
	RedisConn            *redis.Client //链接
	Proxy2layerQueueName string        //队列名称
	Layer2proxyQueueName string        //队列名称
	IdBlackListHash      string        //用户黑名单hash表
	IpBlackListHash      string        //IP黑名单Hash表
	IdBlackListQueue     string        //用户黑名单队列
	IpBlackListQueue     string        //IP黑名单队列
}


type SecKillConf struct {
	CookieSecretKey string

	ReferWhiteList []string //白名单


	AccessLimitConf AccessLimitConf

	RWBlackLock                  sync.RWMutex
	WriteProxy2LayerGoroutineNum int
	ReadProxy2LayerGoroutineNum  int

}

//访问限制
type AccessLimitConf struct {
	IPSecAccessLimit   int //IP每秒钟访问限制
	UserSecAccessLimit int //用户每秒钟访问限制
	IPMinAccessLimit   int //IP每分钟访问限制
	UserMinAccessLimit int //用户每分钟访问限制
}