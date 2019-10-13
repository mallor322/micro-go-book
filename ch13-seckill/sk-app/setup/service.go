package setup

import (
	"github.com/gin-gonic/gin"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/bootstrap"
	"log"
)

//初始化Http服务
func InitServer() {
	router := gin.Default()
	setupRouter(router)
	err := router.Run(bootstrap.HttpConfig.Host + bootstrap.HttpConfig.Port)
	if err != nil {
		log.Printf("Init http server. Error : %v", err)
	}
}
