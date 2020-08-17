module github.com/sorborail/shippy/shippy-cli-consignment

go 1.15

require (
	github.com/micro/go-micro/v2 v2.9.1
)

replace (
	google.golang.org/grpc => google.golang.org/grpc v1.27.0
	github.com/coreos/etcd => github.com/ozonru/etcd v3.3.20-grpc1.27-origmodule+incompatible
)
