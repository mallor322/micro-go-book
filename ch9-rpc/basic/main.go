package main

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch9-rpc/basic/server"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	arith := new(server.Library)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &server.BookInfoParams{1}
	var reply *server.BookInfo
	err = client.Call("Library.GetBookInfo", args, reply)
	if err != nil {
		log.Fatal("Library error:", err)
	}
	fmt.Printf("Library get Book : bookId is %d and name is %s", reply.BookId, reply.Name)
}
