package main

import (
	"fmt"
	register "github.com/longjoy/micro-go-book/ch13-seckill/pkg/discover"
	"github.com/longjoy/micro-go-book/ch13-seckill/sk-core/setup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	errChan := make(chan error)

	setup.InitZk()
	setup.InitRedis()
	setup.RunService()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	//服务退出取消注册
	register.Deregister()
	fmt.Println(error)
}
