package config

import (
	"github.com/gohouse/gorose/v2"
	"go.etcd.io/etcd/clientv3"
)

var SecAdminConfCtx = &SecAdminConf{}

type SecAdminConf struct {
	DbConf   *DbConf
	EtcdConf *EtcdConf
}

//数据库配置
type DbConf struct {
	DbConn gorose.Connection //链接
}

//Etcd配置
type EtcdConf struct {
	EtcdConn          *clientv3.Client //链接
	EtcdSecProductKey string           //商品键
}
