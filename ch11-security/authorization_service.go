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
	//http.HandleFunc("/login", login)
	//http.HandleFunc("/oauth/authorize", authorize)
	http.Handle("/simple", OAuthContextMiddleware(oauthAuthorizationMiddleware(http.HandlerFunc(simple))))
	http.Handle("/admin", OAuthContextMiddleware(oauthAuthorizationMiddleware(adminAuthorityMiddleware(http.HandlerFunc(admin)))))

	http.Handle("/oauth/token", OAuthContextMiddleware(clientAuthorizationMiddleware(http.HandlerFunc(getOAuthToken))))
	http.Handle("/oauth/check_token", OAuthContextMiddleware(clientAuthorizationMiddleware(http.HandlerFunc(checkToken))))

	err := basic.Server.ListenAndServe()
	if err != nil{
		basic.Logger.Println("Service is going to close...")
	}
}


func simple(writer http.ResponseWriter, reader *http.Request)  {
	writer.Write([]byte("simple data"))
}

func admin(writer http.ResponseWriter, reader *http.Request)  {
	writer.Write([]byte("admin data"))
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
			return
		}

		decodeBytes, err := base64.StdEncoding.DecodeString(authorization)

		if err != nil{
			http.Error(w, "Please provide correct base64 encoding", http.StatusForbidden)
			return
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
			return
		}
		w.(OAuthContextResponseWriter).Set("client", clientDetails)
		next.ServeHTTP(w, r)
	})
}


func oauthAuthorizationMiddleware(next http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		basic.Logger.Println("Executing authorization handler")

		// 获取访问令牌
		accessTokenValue := r.Header.Get("Authorization")

		if accessTokenValue == ""{
			http.Error(w, "Please provide access token in authorization", http.StatusForbidden)
			return
		}

		// 获取令牌对应的用户信息和客户端信息
		oauth2Details, err := tokenService.GetOAuth2DetailsByAccessToken(accessTokenValue)

		if err != nil{
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}

		w.(OAuthContextResponseWriter).Set("oauth2Details", oauth2Details)
		next.ServeHTTP(w, r)
	})
}

func adminAuthorityMiddleware(next http.Handler) http.Handler{

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 获取用户信息和客户端信息
		oauth2Details := w.(OAuthContextResponseWriter).Value("oauth2Details")

		if oauth2Details == nil{
			http.Error(w, "Please provide access token in authorization", http.StatusForbidden)
			return
		}

		userDetails := oauth2Details.(*OAuth2Details).User
		for _, value := range userDetails.Authorities{
			if value == "Admin"{
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Authority is not enough", http.StatusForbidden)

	})
}


func getOAuthToken(writer http.ResponseWriter, reader *http.Request)  {

	clientDetails := writer.(OAuthContextResponseWriter).Value("client")

	grantType := reader.URL.Query().Get("grant_type")
	if grantType == ""{
		writer.Write([]byte("Please Input grant_type"))
		return
	}

	for _, v := range clientDetails.(*ClientDetails).AuthorizedGrantTypes{
		if v == grantType{
			oauthToken, err := tokenGrant.grant(grantType, clientDetails.(*ClientDetails), reader)
			if err == nil{
				writer.Header().Set("Content-Type", "application/json")
				json.NewEncoder(writer).Encode(oauthToken)
			}else {
				writer.Write([]byte(err.Error()))
			}
			return

		}
	}
	writer.WriteHeader(403)
	writer.Write([]byte("ClientId " + clientDetails.(*ClientDetails).ClientId + " No Support " + grantType))

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
	// 内置一个客户端
	clientDetailsService = NewInMemoryClientDetailService([] *ClientDetails{&ClientDetails{
		"clientId",
		"clientSecret",
		1800,
		18000,
		"http://127.0.0.1",
		[] string{"password", "refresh_token"},
	}})
	// 内置两个用户
	userDetailsService = NewInMemoryUserDetailService([] *UserDetails{{
		Username:    "simple",
		Password:    "123456",
		UserId:      1,
		Authorities: []string{"Simple"},
	},
		{
			Username:    "admin",
			Password:    "123456",
			UserId:      1,
			Authorities: []string{"Admin"},
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
	// 定义 password 和 refresh_token 令牌生成器
	tokenGrantDict["password"] = &UsernamePasswordTokenGranter{
		supportGrantType:"password",
		userDetailsService:userDetailsService,
		tokenService:tokenService,
	}
	tokenGrantDict["refresh_token"] = &RefreshTokenGranter{
		supportGrantType:"refresh_token",
		tokenService:tokenService,
	}
	tokenGrant = NewComposeTokenGranter(tokenGrantDict)
	basic.StartService("Authorization", "127.0.0.1", 10098, startAuthorizationHttpListener)
}