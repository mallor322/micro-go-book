package main

import (
	"encoding/base64"
	"github.com/keets2012/Micro-Go-Pracrise/basic"
	"net/http"
	"strconv"
	"strings"
)

func startAuthorizationHttpListener(host string, port int)  {
	basic.Server = &http.Server{
		Addr: host + ":" +strconv.Itoa(port),
	}
	http.HandleFunc("/health", basic.CheckHealth)
	http.HandleFunc("/discovery", basic.DiscoveryService)
	http.Handle("/oauth/token", OAuthContextMiddleware(clientAuthorizationMiddleware(http.HandlerFunc(getOAuthToken))))
	err := basic.Server.ListenAndServe()
	if err != nil{
		basic.Logger.Println("Service is going to close...")
	}
}



func clientAuthorizationMiddleware(next http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basic.Logger.Println("Executing client authorization handler")

		// 对客户端信息进行校验
		// clientId:clientSecret Base64加密

		authorization := r.Header.Get("Authorization")

		if authorization == ""{
			http.Error(w, "Please provide the clientId and clientSecret in authorization", http.StatusForbidden)
		}

		decodeBytes, err := base64.StdEncoding.DecodeString(authorization)

		if err != nil{
			http.Error(w, "Please provide correct base64 encoding", http.StatusForbidden)
		}


		decodeStrings := strings.SplitN(string(decodeBytes), ":", 2)

		clientId := decodeStrings[0]
		clientSecret := decodeStrings[1]

		clientDetails, err := clientDetailsService.GetClientDetailByClientId(clientId)

		if err != nil{
			http.Error(w, err.Error(), http.StatusForbidden)
		}

		if !clientDetails.IsMatch(clientId, clientSecret){
			http.Error(w, "Please provide correct client information", http.StatusForbidden)

		}
		next.ServeHTTP(w, r)

	})
	
}


func getOAuthToken(writer http.ResponseWriter, reader *http.Request)  {

	return

	
}



var clientDetailsService ClientDetailService
var userDetailsService UserDetailsService


func main()  {
	clientDetailsService = NewInMemoryClientDetailService([] *ClientDetails{&ClientDetails{
		"clientId",
			"clientSercet",
			1800,
			18000,
			"http://127.0.0.1",
			[] string{"password", "authorization_code"},
	}})
	userDetailsService = NewInMemoryUserDetailService([] *UserDetails{&UserDetails{
		Username:"xuan",
		Password:"123456",
		UserId:1,
	}})
	basic.StartService("Authorization", "127.0.0.1", 10087, startAuthorizationHttpListener)
}
