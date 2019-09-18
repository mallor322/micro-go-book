package main

import (
	"net/http"
)

type OAuthContextResponseWriter struct {
	http.ResponseWriter
	context map[string] interface{}
}

func (w OAuthContextResponseWriter) Value(key string) interface{} {
	return w.context[key]
}

func (w OAuthContextResponseWriter) Set(key string, value interface{})  {
	w.context[key] = value
}


func OAuthContextMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		oauthContextResponseWriter := OAuthContextResponseWriter{
			ResponseWriter: w,
			context:        make(map[string]interface{}),
		}
		next.ServeHTTP(oauthContextResponseWriter, r)
	})
}
