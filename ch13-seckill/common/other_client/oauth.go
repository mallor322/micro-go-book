package main

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
)

func main() {

	client, _ := NewOAuthClient("oauth", nil)
	resp, err := client.CheckToken(&pb.CheckTokenRequest{
		Token: "OOOOOKk",
	})
	if err == nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(resp)
	}

}
