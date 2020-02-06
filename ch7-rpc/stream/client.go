package stream

import (
	"context"
	"fmt"
	"github.com/longjoy/micro-go-book/ch7-rpc/stream-pb"
	"google.golang.org/grpc"
	"io"
	"log"
	"strconv"
)

func main() {
	serviceAddress := "127.0.0.1:50052"
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		panic("connect error")
	}
	defer conn.Close()
	stringClient := stream_pb.NewStringServiceClient(conn)
	stringReq := &stream_pb.StringRequest{A: "A", B: "B"}
	stream, _ := stringClient.LotsOfServerStream(context.Background(), stringReq)
	for {
		item, stream_error := stream.Recv()
		if stream_error == io.EOF {
			break
		}
		if stream_error != nil {
			log.Printf("failed to recv: %v", stream_error)
		}
		fmt.Printf("StringService Concat : %d concat %d = %d", stringReq.A, stringReq.B, item.ProtoMessage)
	}

}

func sendClientStreamRequest(client stream_pb.StringServiceClient) {
	stream, err := client.LotsOfClientStream(context.Background())
	for i := 0; i < 100; i++ {
		if err != nil {
			log.Printf("failed to call: %v", err)
			break
		}
		stream.Send(&stream_pb.StringRequest{A: strconv.Itoa(i), B: strconv.Itoa(i + 1)})
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Printf("failed to recv: %v", err)
	}
	log.Printf("ret is : %s", reply.Ret)
}

func sendClientAndServerStreamRequest(client stream_pb.StringServiceClient) {
	var err error
	stream, err := client.LotsOfServerAndClientStream(context.Background())
	if err != nil {
		log.Printf("failed to call: %v", err)
		return
	}
	var i int64
	for {
		err1 := stream.Send(&stream_pb.StringRequest{A: strconv.Itoa(i), B: strconv.Itoa(i + 1)})
		if err1 != nil {
			log.Printf("failed to send: %v", err)
			break
		}
		reply, err2 := stream.Recv()
		if err2 != nil {
			log.Printf("failed to recv: %v", err)
			break
		}
		log.Printf("Ret is : %s", reply.Ret)
		i++
	}
}
