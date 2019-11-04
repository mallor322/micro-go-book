package grpc

import (
	"context"
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch9-rpc/pb"
	"google.golang.org/grpc"
)

func main() {
	serviceAddress := "127.0.0.1:50052"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	bookClient := pb.NewStringServiceClient(conn)
	stringReq := &pb.StringRequest{A: "A", B: "B"}
	reply, _ := bookClient.Concat(context.Background(), stringReq)
	fmt.Printf("StringService Concat : %d concat %d = %d", stringReq.A, stringReq.B, reply.Ret)
}
