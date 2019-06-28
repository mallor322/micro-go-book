module github.com/keets2012/Micro-Go-Pracrise

go 1.12

require (
	github.com/go-kit/kit v0.9.0
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/gorilla/mux v1.6.2
	github.com/hashicorp/consul/api v1.1.0
	github.com/juju/ratelimit v1.0.1
	github.com/openzipkin/zipkin-go v0.1.6
	github.com/pborman/uuid v1.2.0
	github.com/prometheus/client_golang v1.0.0
	golang.org/x/time v0.0.0-00010101000000-000000000000
)

replace golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
