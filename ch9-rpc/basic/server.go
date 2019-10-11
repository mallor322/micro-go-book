package main

import (
	"github.com/keets2012/Micro-Go-Pracrise/ch9-rpc/basic/server"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	arith := new(server.Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}
