module cloud-store.com

go 1.13

require (
	github.com/aliyun/aliyun-oss-go-sdk v0.0.0-20190307165228-86c17b95fcd5
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/contrib v0.0.0-20201101042839-6a891bf89f19
	github.com/gin-gonic/gin v1.6.3
	github.com/go-bindata/go-bindata v3.1.2+incompatible // indirect
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.2 // indirect
	github.com/json-iterator/go v1.1.9
	github.com/jteeuwen/go-bindata v3.0.7+incompatible // indirect
	github.com/juju/ratelimit v1.0.1
	github.com/micro/cli v0.2.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins/registry/consul v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/breaker/hystrix v0.0.0-20200119172437-4fe21aa238fd
	github.com/micro/go-plugins/wrapper/ratelimiter/ratelimit v0.0.0-20200119172437-4fe21aa238fd
	github.com/mitchellh/mapstructure v1.1.2
	github.com/moxiaomomo/go-bindata-assetfs v1.0.0
	github.com/olivere/elastic v6.2.35+incompatible
	github.com/sqs/goreturns v0.0.0-20181028201513-538ac6014518 // indirect
	github.com/streadway/amqp v1.0.0
	go.uber.org/multierr v1.6.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/tools v0.0.0-20201105220310-78b158585360 // indirect
	golang.org/x/tools/gopls v0.5.2 // indirect
	google.golang.org/grpc v1.33.2 // indirect
	google.golang.org/grpc/examples v0.0.0-20201106192519-9c2f82d9a79c // indirect
	google.golang.org/protobuf v1.25.0
	gopkg.in/amz.v1 v1.0.0-20150111123259-ad23e96a31d2
	gopkg.in/olivere/elastic.v5 v5.0.86
	honnef.co/go/tools v0.0.1-2020.1.6 // indirect
	mvdan.cc/gofumpt v0.0.0-20201107090320-a024667a00f1 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
