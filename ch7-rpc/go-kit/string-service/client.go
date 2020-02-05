package main

import (
	"context"
	"flag"
	"fmt"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/longjoy/micro-go-book/ch13-seckill/pb"
	"google.golang.org/grpc"
	"time"
)

func main() {
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

	svr := NewStringClient(conn)
	result, err := svr.Concat(ctx, "A", "B")
	if err != nil {
		fmt.Println("Check error", err.Error())

	}

	fmt.Println("result=", result)
}

func NewStringClient(conn *grpc.ClientConn) Service {

	var ep = grpctransport.NewClient(conn,
		"pb.StringService",
		"Concat",
		DecodeStringRequest,
		EncodeStringResponse,
		pb.UserResponse{},
	).Endpoint()

	userEp := StringEndpoints{
		StringEndpoint: ep,
	}
	return userEp
}
