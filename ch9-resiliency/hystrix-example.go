package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"time"
)

func main(){

	err := hystrix.Do("test_command", func() error {
		// 远程调用&或者其他需要保护的方法
		return nil
	}, func(err error) error{
		// 失败回滚方法
		return nil
	})


	resultChan := make(chan interface{}, 1)
	errChan := hystrix.Go("test_command", func() error {
		// 远程调用&或者其他需要保护的方法
		resultChan <- "success"
		return nil
	}, func(e error) error {
		// 失败回滚方法
		return nil
	})

	select {
	case err := <- errChan:
		// 执行失败
	case result := <- resultChan :
		// 执行成功
	case <-time.After(2 * time.Second): // 超时设置
		fmt.Println("Time out")
		return
	}




}