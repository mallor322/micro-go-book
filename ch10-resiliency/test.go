package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

func main()  {


	// 注册watch，监控所有服务的变化


	var params map[string]interface{}
	json.Unmarshal([]byte(`{"type":"service", "service" : "OO"}`), &params)
	plan, _ := watch.Parse(params)
	plan.Handler = func(u uint64, i interface{}) {
		if i == nil{
			return
		}

		v, ok := i.([]*api.ServiceEntry)
		if !ok || len(v) == 0 {
			return // ignore
		}

		fmt.Println(len(v))


	}
	plan.Run("127.0.0.1:8500")
	fmt.Println("OK")



}
