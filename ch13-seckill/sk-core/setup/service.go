package setup

import (
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/sk-core/service/srv_redis"
)

func RunService() {
	//启动处理线程
	srv_redis.RunProcess()

	//ctx := context.Background()
	//errChan := make(chan error)

	//var svc Service
	//svc = ArithmeticService{}
	//endpoint := MakeArithmeticEndpoint(svc)

	//var logger log.Logger
	//{
	//	logger = log.NewLogfmtLogger(os.Stderr)
	//	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	//	logger = log.With(logger, "caller", log.DefaultCaller)
	//}
	//
	//r := MakeHttpHandler(ctx, endpoint, logger)

	//go func() {
	//	fmt.Println("Http Server start at port:9000")
	//	handler := r
	//	errChan <- http.ListenAndServe(":9000", handler)
	//}()

	//go func() {
	//	c := make(chan os.Signal, 1)
	//	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	//	errChan <- fmt.Errorf("%s", <-c)
	//}()
	//
	//fmt.Println(<-errChan)
}
