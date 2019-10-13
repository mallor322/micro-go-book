package setup

import (
	"github.com/go-redis/redis"
	conf "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/config"
	"log"
)

//初始化redis
func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     "39.98.179.73:6379", //conf.Redis.Host,
		Password: "082203",            //conf.Redis.Password,
		DB:       conf.Redis.Db,
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Printf("Connect redis failed. Error : %v", err)
	}
	conf.Redis.RedisConn = client
}
