package grpc

import (
	"flag"
	"fmt"
	string_service "github.com/longjoy/micro-go-book/ch7-rpc/grpc/string-service"
	"github.com/longjoy/micro-go-book/ch7-rpc/pb"
	"github.com/prometheus/common/log"
	"google.golang.org/grpc"
	"net"
)

func main() {
	flag.Parse()
	port := ":8081"
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	stringService := new(string_service.StringServer)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
