package stream

import (
	"flag"
	"fmt"
	pb "github.com/longjoy/micro-go-book/ch7-rpc/stream-pb"
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
	stringService := new(StringServer)
	pb.RegisterStringServiceServer(grpcServer, stringService)
	grpcServer.Serve(lis)
}
