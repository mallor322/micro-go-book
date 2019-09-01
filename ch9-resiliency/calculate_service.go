package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func startCalculateHttpListener(host string, port int)  {
	server = &http.Server{
		// GetLocalIpAddress用于获取本地IP，可以手动写入
		Addr: host + ":" +strconv.Itoa(port),
	}
	http.HandleFunc("/health", checkHealth)
	http.HandleFunc("/calculate", calculate)
	http.HandleFunc("/discovery", discoveryService)
	err := server.ListenAndServe()
	if err != nil{
		logger.Println("Service is going to close...")
	}
}


func calculate(writer http.ResponseWriter, reader *http.Request)  {
	a, _:= strconv.Atoi(reader.URL.Query().Get("a"))
	b, _:= strconv.Atoi(reader.URL.Query().Get("b"))
	_, err := fmt.Fprintln(writer, a + b)
	if err != nil{
		logger.Println(err)
	}
}


func main()  {

	startService("Calculate", "127.0.0.1", 10085, startCalculateHttpListener)

}
