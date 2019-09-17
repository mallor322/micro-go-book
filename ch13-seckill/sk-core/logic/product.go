package logic

import (
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-core/config"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-core/service/srv_limit"
	"context"
	"encoding/json"
	"log"
	"time"
)
//"go.etcd.io/etcd/mvcc/mvccpb"

//从Etcd中加载商品数据
func LoadProductFromEtcd() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	//从etcd获取商品数据
	rsp, err := config.SecLayerCtx.EtcdConf.EtcdConn.Get(ctx, config.SecLayerCtx.EtcdConf.EtcdSecProductKey)
	if err != nil {
		log.Printf("get [%s] from etcd failed. Error : %v", config.SecLayerCtx.EtcdConf.EtcdSecProductKey, err)
		return
	}

	//结构转换
	var secProductInfo []*config.SecProductInfoConf
	for _, v := range rsp.Kvs {
		err = json.Unmarshal(v.Value, &secProductInfo)
		if err != nil {
			log.Printf("Unmarshal sec product info failed. Error : %v", err)
			return
		}
	}

	updateSecProductInfo(secProductInfo)
	log.Printf("update product info success, data : %v", secProductInfo)

	initSecProductWatcher()
	log.Printf("init ectd watcher success.")

	return
}

//更新商品信息
func updateSecProductInfo(secProductInfo []*config.SecProductInfoConf) {
	tmp := make(map[int]*config.SecProductInfoConf, 1024)

	for _, v := range secProductInfo {
		productInfo := v
		productInfo.SecLimit = &srv_limit.SecLimit{}
		tmp[v.ProductId] = productInfo
	}

	config.SecLayerCtx.RWSecProductLock.Lock()
	config.SecLayerCtx.SecLayerConf.SecProductInfoMap = tmp
	config.SecLayerCtx.RWSecProductLock.Unlock()
}

//监控商品变化
func initSecProductWatcher() {
	go watchSecProductKey()
}

func watchSecProductKey() {
	key := config.SecLayerCtx.EtcdConf.EtcdSecProductKey

	//var err error
	for {
		rch := config.SecLayerCtx.EtcdConf.EtcdConn.Watch(context.Background(), key)
		var secProductInfo []*config.SecProductInfoConf
		var getConfSucc = true

		for wrsp := range rch {
			for _, ev := range wrsp.Events {
				//删除事件
				//if ev.Type.EnumDescriptor() == mvccpb.DELETE.EnumDescriptor() {
				//	log.Printf("key[%s] 's config deleted", key)
				//	continue
				//}
				//
				////更新事件
				//if ev.Type.EnumDescriptor() == mvccpb.PUT.EnumDescriptor() && string(ev.Kv.Key) == key {
				//	err = json.Unmarshal(ev.Kv.Value, &secProductInfo)
				//	if err != nil {
				//		log.Printf("key [%s], Unmarshal[%s]. Error : %v", key, err)
				//		getConfSucc = false
				//		continue
				//	}
				//}
				log.Printf("get config from etcd, %s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			}

			if getConfSucc {
				log.Printf("get config from etcd success, %v", secProductInfo)
				updateSecProductInfo(secProductInfo)
			}
		}
	}
}
