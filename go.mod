module github.com/keets2012/Micro-Go-Pracrise

go 1.12

require (
	github.com/afex/hystrix-go v0.0.0-20180502004556-fa1af6a1f4f5
	github.com/apache/thrift v0.12.0 // indirect
	github.com/go-kit/kit v0.9.0
	github.com/go-logfmt/logfmt v0.4.0 // indirect
	github.com/go-redis/redis v6.15.5+incompatible
	github.com/golang/protobuf v1.3.1
	github.com/gorilla/mux v1.6.2
	github.com/hashicorp/consul/api v1.1.0
	github.com/juju/ratelimit v1.0.1
	github.com/opentracing-contrib/go-observer v0.0.0-20170622124052-a52f23424492 // indirect
	github.com/opentracing/opentracing-go v1.1.0
	github.com/openzipkin-contrib/zipkin-go-opentracing v0.3.5
	github.com/openzipkin/zipkin-go v0.1.6
	github.com/pborman/uuid v1.2.0
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.0.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/viper v1.4.0
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	go.etcd.io/etcd v3.3.15+incompatible
	golang.org/x/time v0.0.0-20190308202827-9d24e82272b4
)

replace golang.org/x/time => github.com/golang/time v0.0.0-20190308202827-9d24e82272b4
