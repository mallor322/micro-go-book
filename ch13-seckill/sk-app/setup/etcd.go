package setup

import (
	"context"
	"encoding/json"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	conf "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/config"
	"log"
	"time"
)

//初始化Etcd
func InitEtcd() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"39.98.179.73:2379"}, // conf.Etcd.Host
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("Connect etcd failed. Error : %v", err)
	}
	log.Print("Connect etcd sucess")
	conf.Etcd.EtcdConn = cli
	loadSecConf(cli)
	go waitSecProductKey(cli, conf.Etcd.EtcdSecProductKey)
}

//加载秒杀商品信息
func loadSecConf(cli *clientv3.Client) {
	log.Printf("Connect etcd sucess %s", conf.Etcd.EtcdSecProductKey)
	rsp, err := cli.Get(context.Background(), "sec_kill_product") //conf.Etcd.EtcdSecProductKey
	log.Print("Connect etcd sucess")
	if err != nil {
		log.Print("Connect etcd sucess")
		log.Printf("get product info failed, err : %v", err)
		return
	}
	log.Printf("get product info ")
	var secProductInfo []*conf.SecProductInfoConf
	for _, v := range rsp.Kvs {
		err := json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			log.Printf("unmarshal json failed, err : %v", err)
			return
		}
	}
	log.Print("loadSecConf finished")

	updateSecProductInfo(secProductInfo)
}

//监听秒杀商品配置
func waitSecProductKey(cli *clientv3.Client, key string) {
	for {
		rch := cli.Watch(context.Background(), key)
		var secProductInfo []*conf.SecProductInfoConf
		var getConfSucc = true

		for wrsp := range rch {
			for _, ev := range wrsp.Events {
				//删除事件
				if ev.Type == mvccpb.DELETE {
					continue
				}

				//更新事件
				if ev.Type == mvccpb.PUT && string(ev.Kv.Key) == key {
					err := json.Unmarshal(ev.Kv.Value, &secProductInfo)
					if err != nil {
						getConfSucc = false
						continue
					}
				}
			}

			if getConfSucc {
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}

//更新秒杀商品信息
func updateSecProductInfo(secProductInfo []*conf.SecProductInfoConf) {
	tmp := make(map[int]*conf.SecProductInfoConf, 1024)
	for _, v := range secProductInfo {
		tmp[v.ProductId] = v
	}
	conf.SecKill.RWBlackLock.Lock()
	conf.SecKill.SecProductInfoMap = tmp
	conf.SecKill.RWBlackLock.Unlock()
}
