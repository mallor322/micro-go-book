package main

import (
	"fmt"
	"github.com/keets2012/Micro-Go-Pracrise/basic"
	"net/http"
	"strconv"
)

func startCalculateHttpListener(host string, port int) {
	basic.Server = &http.Server{
		Addr: host + ":" + strconv.Itoa(port),
	}
	http.HandleFunc("/health", basic.CheckHealth)
	http.HandleFunc("/calculate", calculate)
	http.HandleFunc("/discovery", basic.DiscoveryService)
	err := basic.Server.ListenAndServe()
	if err != nil {
		basic.Logger.Println("Service is going to close...")
	}
}

func calculate(writer http.ResponseWriter, reader *http.Request) {
	a, _ := strconv.Atoi(reader.URL.Query().Get("a"))
	b, _ := strconv.Atoi(reader.URL.Query().Get("b"))
	_, err := fmt.Fprintln(writer, a+b)
	if err != nil {
		basic.Logger.Println(err)
	}
}

func main() {
	basic.StartService("Calculate", "127.0.0.1", 10085, startCalculateHttpListener)
}
