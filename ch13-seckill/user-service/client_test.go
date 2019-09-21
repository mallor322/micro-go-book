package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/client"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var (
		grpcAddr = flag.String("addr", ":9008", "gRPC address")
	)
	flag.Parse()

	ctx := context.Background()

	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	svr := client.UserCheck(conn)
	result, err := svr.Check(ctx, "Add", "pps")
	if err != nil {
		fmt.Println("Check error", err.Error())

	}

	fmt.Println("result=", result)
}
