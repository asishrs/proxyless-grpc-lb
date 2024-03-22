module github.com/asishrs/proxyless-grpc-lb/xds-server

go 1.15

require (
	github.com/asishrs/proxyless-grpc-lb/common v0.0.0-20200824005131-5718db7e19f5
	github.com/envoyproxy/go-control-plane v0.11.1-0.20230524094728-9239064ad72f
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.5.3
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.3.0
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.17.0 // indirect
	google.golang.org/grpc v1.56.3
	gopkg.in/ini.v1 v1.60.1 // indirect
	k8s.io/apimachinery v0.18.8
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20200821003339-5e75c0163111 // indirect
)

replace (
	github.com/envoyproxy/go-control-plane => github.com/asishrs/go-control-plane v0.9.7
	github.com/googleapis/gnostic => github.com/googleapis/gnostic v0.4.0
	k8s.io/apimachinery => k8s.io/apimachinery v0.18.0
	k8s.io/client-go => k8s.io/client-go v0.18.0
)
