package mysql

import (
	"github.com/gohouse/gorose"
)

var DbConfCtx = &DbConf{}

//数据库配置
type DbConf struct {
	DbConn gorose.Connection //链接
}
