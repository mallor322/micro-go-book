package client

import (
	"context"
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/pb"
	"testing"
)

func TestUserClientImpl_CheckUser(t *testing.T) {



	client, _ := NewUserClient("user", nil)

	if response, err := client.CheckUser(context.Background(), &pb.UserRequest{
		Username:"xuan",
		Password:"xuan",
	});  err == nil {
		fmt.Println(response.Result)
	}else {
		fmt.Println(err.Error())
	}


}
