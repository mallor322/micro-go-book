package main

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/ch9-rpc/basic/string-service"
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	stringReq := &service.StringRequest{"A", "B"}
	// Synchronous call
	var reply string
	err = client.Call("StringService.Concat", stringReq, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("StringService Concat : %d concat %d = %d", stringReq.A, stringReq.B, reply)

	stringReq = &service.StringRequest{"ACD", "BDF"}
	call := client.Go("StringService.Concat", stringReq, &reply, nil)
	_ := <-call.Done
	fmt.Printf("StringService Diff : %d diff %d = %d", stringReq.A, stringReq.B, reply)

}
