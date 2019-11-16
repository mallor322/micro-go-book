package main

import (
	"context"
	"flag"
	"fmt"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/client"
	localconfig "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/config"
	"github.com/openzipkin/zipkin-go"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	var (
		grpcAddr = flag.String("addr", ":9008", "gRPC address")
	)
	flag.Parse()
	tr := localconfig.ZipkinTracer
	parentSpan := tr.StartSpan("test")
	ctx := zipkin.NewContext(context.Background(), parentSpan)

	clientTracer := kitzipkin.GRPCClientTrace(tr, kitzipkin.Name("grpc-transport"))
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	svr := client.UserCheck(conn, clientTracer)

	result, err := svr.Check(ctx, "Add", "pps")
	if err != nil {
		fmt.Println("Check error", err.Error())

	}

	fmt.Println("result=", result)
}
