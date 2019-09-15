### 秒杀性能配置

```
serviceMap := viper.GetStringMap("service")
			ipSecAccessLimit, _ := serviceMap["ip_sec_access_limit"].(int64)
			ipMinAccessLimit, _ := serviceMap["ip_min_access_limit"].(int64)
			userSecAccessLimit, _ := serviceMap["user_sec_access_limit"].(int64)
			userMinAccessLimit, _ := serviceMap["user_min_access_limit"].(int64)
			writeProxy2layerGoroutineNum, _ := serviceMap["write_proxy2layer_goroutine_num"].(int64)
			readProxy2layerGoroutineNum, _ := serviceMap["read_proxy2layer_goroutine_num"].(int64)
			cookieSecretKey, _ := serviceMap["cookie_secretkey"].(string)
			referWhitelist, _ := serviceMap["refer_whitelist"].(string)
			setup.InitServiceConfig(ipSecAccessLimit, ipMinAccessLimit, userSecAccessLimit, userMinAccessLimit,
				writeProxy2layerGoroutineNum, readProxy2layerGoroutineNum, cookieSecretKey, referWhitelist)

			redisMap := viper.GetStringMap("redis")
			hostRedis, _ := redisMap["host"].(string)
			pwdRedis, _ := redisMap["password"].(string)
			dbRedis, _ := redisMap["db"].(int)
			proxy2layerQueueNameRedis, _ := redisMap["proxy2layer_queue_name"].(string)
			layer2proxyQueueNameRedis, _ := redisMap["layer2proxy_queue_name"].(string)
			idBlackListHashRedis, _ := redisMap["id_black_list_hash"].(string)
			ipBlackListHashRedis, _ := redisMap["ip_black_list_hash"].(string)
			idBlackListQueueRedis, _ := redisMap["id_black_list_queue"].(string)
			ipBlackListQueueRedis, _ := redisMap["ip_black_list_queue"].(string)
			setup.InitRedis(hostRedis, pwdRedis, dbRedis, proxy2layerQueueNameRedis, layer2proxyQueueNameRedis,
				idBlackListHashRedis, ipBlackListHashRedis, idBlackListQueueRedis, ipBlackListQueueRedis)

			etcdMap := viper.GetStringMap("etcd")
			hostEtcd, _ := etcdMap["host"].(string)
			productKey, _ := etcdMap["product_key"].(string)
			setup.InitEtcd(hostEtcd, productKey)

			httpMap := viper.GetStringMap("http")
			hostHttp, _ := httpMap["host"].(string)
			setup.InitServer(hostHttp)
		},
	}
```

```
var RootCmd = &cobra.Command{
		Use:   "SKLayer Server",
		Short: "SKLayer Server",
		Long:  "SKLayer Server",
		Run: func(cmd *cobra.Command, args []string) {
			serviceMap := viper.GetStringMap("service")
			writeProxy2layerGoroutineNum, _ := serviceMap["write_proxy2layer_goroutine_num"].(int64)
			readLayer2proxyGoroutineNum, _ := serviceMap["read_layer2proxy_goroutine_num"].(int64)
			handleUserGoroutineNum, _ := serviceMap["handle_user_goroutine_num"].(int64)
			handle2WriteChanSize, _ := serviceMap["handle2write_chan_size"].(int64)
			maxRequestWaitTimeout, _ := serviceMap["max_request_wait_timeout"].(int64)
			sendToWriteChanTimeout, _ := serviceMap["send_to_write_chan_timeout"].(int64)
			sendToHandleChanTimeout, _ := serviceMap["send_to_handle_chan_timeout"].(int64)
			secKillTokenPassWd, _ := serviceMap["seckill_token_passwd"].(string)
			setup.InitService(writeProxy2layerGoroutineNum, readLayer2proxyGoroutineNum, handleUserGoroutineNum,
				handle2WriteChanSize, maxRequestWaitTimeout, sendToWriteChanTimeout, sendToHandleChanTimeout, secKillTokenPassWd)

			redisMap := viper.GetStringMap("redis")
			hostRedis, _ := redisMap["host"].(string)
			pwdRedis, _ := redisMap["password"].(string)
			dbRedis, _ := redisMap["db"].(int)
			proxy2layerQueueName, _ := redisMap["proxy2layer_queue_name"].(string)
			layer2proxyQueueName, _ := redisMap["layer2proxy_queue_name"].(string)
			setup.InitRedis(hostRedis, pwdRedis, dbRedis, proxy2layerQueueName, layer2proxyQueueName)

			etcdMap := viper.GetStringMap("etcd")
			hostEtcd, _ := etcdMap["host"].(string)
			productKey, _ := etcdMap["product_key"].(string)
			setup.InitEtcd(hostEtcd, productKey)

			setup.RunService()
		},
	}
```