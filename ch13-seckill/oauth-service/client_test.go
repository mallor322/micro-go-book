package main

import (
	"flag"
	"fmt"
	kitzipkin "github.com/go-kit/kit/tracing/zipkin"
	"github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/common/client"
	localconfig "github.com/keets2012/Micro-Go-Pracrise/ch13-seckill/user-service/config"
	"google.golang.org/grpc"
	"testing"
	"time"
)



func TestOAuthClient(t *testing.T) {
	var (
		grpcAddr = flag.String("addr", ":9008", "gRPC address")
	)
	flag.Parse()
	tr := localconfig.ZipkinTracer

	clientTracer := kitzipkin.GRPCClientTrace(tr, kitzipkin.Name("grpc-transport"))
	conn, err := grpc.Dial(*grpcAddr, grpc.WithInsecure(), grpc.WithTimeout(1*time.Second))
	if err != nil {
		fmt.Println("gRPC dial err:", err)
	}
	defer conn.Close()

	svr := client.CheckToken(conn, clientTracer)
	result, err := svr.GetOAuth2DetailsByAccessToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VyRGV0YWlscyI6eyJVc2VySWQiOjEsIlVzZXJuYW1lIjoidXNlcm5hbWUiLCJQYXNzd29yZCI6IiIsIkF1dGhvcml0aWVzIjpbIkFkbWluIiwiU3VwZXIiXX0sIkNsaWVudERldGFpbHMiOnsiQ2xpZW50SWQiOiJjbGllbnRJZCIsIkNsaWVudFNlY3JldCI6IiIsIkFjY2Vzc1Rva2VuVmFsaWRpdHlTZWNvbmRzIjoxODAwLCJSZWZyZXNoVG9rZW5WYWxpZGl0eVNlY29uZHMiOjE4MDAwLCJSZWdpc3RlcmVkUmVkaXJlY3RVcmkiOiJodHRwOi8vMTI3LjAuMC4xIiwiQXV0aG9yaXplZEdyYW50VHlwZXMiOlsicGFzc3dvcmQiLCJyZWZyZXNoX3Rva2VuIl19LCJSZWZyZXNoVG9rZW4iOnsiUmVmcmVzaFRva2VuIjpudWxsLCJUb2tlblR5cGUiOiJqd3QiLCJUb2tlblZhbHVlIjoiZXlKaGJHY2lPaUpJVXpJMU5pSXNJblI1Y0NJNklrcFhWQ0o5LmV5SlZjMlZ5UkdWMFlXbHNjeUk2ZXlKVmMyVnlTV1FpT2pFc0lsVnpaWEp1WVcxbElqb2lkWE5sY201aGJXVWlMQ0pRWVhOemQyOXlaQ0k2SWlJc0lrRjFkR2h2Y21sMGFXVnpJanBiSWtGa2JXbHVJaXdpVTNWd1pYSWlYWDBzSWtOc2FXVnVkRVJsZEdGcGJITWlPbnNpUTJ4cFpXNTBTV1FpT2lKamJHbGxiblJKWkNJc0lrTnNhV1Z1ZEZObFkzSmxkQ0k2SWlJc0lrRmpZMlZ6YzFSdmEyVnVWbUZzYVdScGRIbFRaV052Ym1Seklqb3hPREF3TENKU1pXWnlaWE5vVkc5clpXNVdZV3hwWkdsMGVWTmxZMjl1WkhNaU9qRTRNREF3TENKU1pXZHBjM1JsY21Wa1VtVmthWEpsWTNSVmNta2lPaUpvZEhSd09pOHZNVEkzTGpBdU1DNHhJaXdpUVhWMGFHOXlhWHBsWkVkeVlXNTBWSGx3WlhNaU9sc2ljR0Z6YzNkdmNtUWlMQ0p5WldaeVpYTm9YM1J2YTJWdUlsMTlMQ0pTWldaeVpYTm9WRzlyWlc0aU9uc2lVbVZtY21WemFGUnZhMlZ1SWpwdWRXeHNMQ0pVYjJ0bGJsUjVjR1VpT2lJaUxDSlViMnRsYmxaaGJIVmxJam9pSWl3aVJYaHdhWEpsYzFScGJXVWlPbTUxYkd4OUxDSmxlSEFpT2pFMU56TXhOREF5TWpnc0ltbHpjeUk2SWxONWMzUmxiU0o5LjE3bGZrZ3RraFBRVTVkYXA5MGtoQjVKUFRFLXU3V0x4aVZrd2FDcG5uLWsiLCJFeHBpcmVzVGltZSI6IjIwMTktMTEtMDdUMjM6MjM6NDguMTEwNDM1KzA4OjAwIn0sImV4cCI6MTU3MzEyNDAyOCwiaXNzIjoiU3lzdGVtIn0.sOgngpp781LNU6JpCxRCcnYTpZ7YfnAr4-aig29JASo")
	if err != nil {
		fmt.Println("Check error", err.Error())
	}else {
		fmt.Println("result=", result.Client.ClientId)
	}
}
