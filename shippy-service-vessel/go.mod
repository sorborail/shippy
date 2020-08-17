module github.com/sorborail/shippy/shippy-service-vessel

go 1.15

require (
	github.com/golang/protobuf v1.4.2
	github.com/micro/go-micro/v2 v2.9.1
	google.golang.org/protobuf v1.25.0
)

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
)
