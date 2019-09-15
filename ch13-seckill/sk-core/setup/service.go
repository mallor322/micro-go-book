package setup

import (
	"SecKill/sk_layer/config"
	"SecKill/sk_layer/service/srv_product"
	"SecKill/sk_layer/service/srv_redis"
	"SecKill/sk_layer/service/srv_user"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"context"
	"os"
	"os/signal"
	"syscall"
)

func InitService(writeProxy2layerGoroutineNum, readLayer2proxyGoroutineNum, handleUserGoroutineNum,
	handle2WriteChanSize, maxRequestWaitTimeout, sendToWriteChanTimeout, sendToHandleChanTimeout int64, secKillTokenPassWd string) {

	config.AppConfig.WriteGoroutineNum = int(writeProxy2layerGoroutineNum)
	config.AppConfig.ReadGoroutineNum = int(readLayer2proxyGoroutineNum)
	config.AppConfig.HandleUserGoroutineNum = int(handleUserGoroutineNum)
	config.AppConfig.Handle2WriteChanSize = int(handle2WriteChanSize)
	config.AppConfig.MaxRequestWaitTimeout = int(maxRequestWaitTimeout)
	config.AppConfig.SendToWriteChanTimeout = int(sendToWriteChanTimeout)
	config.AppConfig.SendToHandleChanTimeout = int(sendToHandleChanTimeout)
	config.AppConfig.TokenPassWd = secKillTokenPassWd

	config.SecLayerCtx.SecLayerConf = config.AppConfig
	config.SecLayerCtx.Read2HandleChan = make(chan *config.SecRequest, config.AppConfig.Read2HandleChanSize)
	config.SecLayerCtx.Handle2WriteChan = make(chan *config.SecResponse, config.AppConfig.Handle2WriteChanSize)
	config.SecLayerCtx.HistoryMap = make(map[int]*srv_user.UserBuyHistory, 10000)
	config.SecLayerCtx.ProductCountMgr = srv_product.NewProductCountMgr()
}

func RunService() {
	//启动处理线程
	srv_redis.RunProcess()


	ctx := context.Background()
	errChan := make(chan error)

	var svc Service
	svc = ArithmeticService{}
	endpoint := MakeArithmeticEndpoint(svc)

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	r := MakeHttpHandler(ctx, endpoint, logger)

	go func() {
		fmt.Println("Http Server start at port:9000")
		handler := r
		errChan <- http.ListenAndServe(":9000", handler)
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(<-errChan)
}
