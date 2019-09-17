package setup

import (
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-app/controller"
	"github.com/gin-gonic/gin"
)

//设置路由
func setupRouter(router *gin.Engine) {
	//秒杀管理
	router.GET("/sec/info", controller.SecInfo)
	router.GET("/sec/list", controller.SecInfoList)
	router.POST("/sec/kill", controller.SecKill)
}
