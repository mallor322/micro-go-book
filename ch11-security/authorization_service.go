package main

import (
	"encoding/base64"
	"encoding/json"
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
	http.HandleFunc("/login", login)
	http.HandleFunc("/oauth/authorize", authorize)
	http.Handle("/oauth/token", OAuthContextMiddleware(clientAuthorizationMiddleware(http.HandlerFunc(getOAuthToken))))
	http.Handle("/oauth/check_token", OAuthContextMiddleware(clientAuthorizationMiddleware(http.HandlerFunc(checkToken))))

	err := basic.Server.ListenAndServe()
	if err != nil{
		basic.Logger.Println("Service is going to close...")
	}
}

func login(writer http.ResponseWriter, reader *http.Request)  {
	username := reader.FormValue("username")
	password := reader.FormValue("password")

	if username == "" || password == ""{
		writer.WriteHeader(403)
	}

	userDetails, err := userDetailsService.GetUserDetailByUsername(username)
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Header().Set("Access-Control-Allow-Headers", "user-cookie")



	if err == nil{
		writer.Write([]byte(strconv.Itoa(userDetails.UserId)))
	}else {
		writer.WriteHeader(403)
	}

}

func authorize(writer http.ResponseWriter, reader *http.Request)  {
	writer.Header().Set("Access-Control-Allow-Origin", "*")
	writer.Write([]byte("www.baidu.com"))
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
		w.(OAuthContextResponseWriter).Set("client", clientDetails)
		next.ServeHTTP(w, r)
	})
}

func getOAuthToken(writer http.ResponseWriter, reader *http.Request)  {

	clientDetails := writer.(OAuthContextResponseWriter).Value("client")

	grantType := reader.URL.Query().Get("grant_type")

	if grantType == ""{
		writer.Write([]byte("Please Input grant_type"))
		return
	}

	oauthToken, err := tokenGrant.grant(grantType, clientDetails.(*ClientDetails), reader)

	if err == nil{
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(oauthToken)
	}else {
		writer.Write([]byte(err.Error()))
	}

}


func checkToken(writer http.ResponseWriter, reader *http.Request) {
	 tokenValue := reader.URL.Query().Get("token")

	 if tokenValue == ""{
		 writer.Write([]byte("Please Input token"))
		 return
	 }

	 oauth2Details, err := tokenService.GetOAuth2DetailsByAccessToken(tokenValue)

	 if err == nil{
		 writer.Header().Set("Content-Type", "application/json")
		 json.NewEncoder(writer).Encode(oauth2Details)
	 }else {
		 writer.Write([]byte(err.Error()))
	 }
}



var clientDetailsService ClientDetailService
var userDetailsService UserDetailsService
var tokenGrant TokenGranter
var tokenService *TokenService


func main()  {
	clientDetailsService = NewInMemoryClientDetailService([] *ClientDetails{&ClientDetails{
		"clientId",
			"clientSecret",
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
	tokenService = &TokenService{
		tokenStore:&JwtTokenStore{
			jwtTokenEnhancer:&JwtTokenEnhancer{
				secretKey: []byte("abcd1234!@#$"),
			},
		},
		tokenEnhancer:&JwtTokenEnhancer{
			secretKey: []byte("abcd1234!@#$"),
		},
	}
	tokenGrantDict := make(map[string] TokenGranter)
	tokenGrantDict["password"] = &UsernamePasswordTokenGranter{
		supportGrantType:"password",
		userDetailsService:userDetailsService,
		tokenService:tokenService,
	}
	tokenGrant = NewComposeTokenGranter(tokenGrantDict)

	basic.StartService("Authorization", "127.0.0.1", 10098, startAuthorizationHttpListener)
}
