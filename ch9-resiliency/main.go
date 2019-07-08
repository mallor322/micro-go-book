package main

import (
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	hystrix.Go("get_baidu", func() error {
		// talk to other services
		res, err := http.Get("https://www.baidu.com/")
		if err != nil {
			fmt.Println("get error")
			return err
		}
		result, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", result)
		return nil
	}, func(err error) error {
		fmt.Println("get an error, handle it")
		return nil
	})

	time.Sleep(2 * time.Second) // 调用Go方法就是起了一个goroutine，这里要sleep一下，不然看不到效果
}
