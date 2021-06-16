module github.com/grobza/proxyless-grpc-lb/xds-server

go 1.16

require (
	github.com/envoyproxy/go-control-plane v0.9.9
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.2.0
	github.com/spf13/viper v1.8.0
	go.uber.org/zap v1.17.0
	google.golang.org/grpc v1.38.0
	k8s.io/apimachinery v0.21.1
	k8s.io/client-go v0.21.1
)
