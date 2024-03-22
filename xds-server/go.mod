module github.com/asishrs/proxyless-grpc-lb/xds-server

go 1.15

require (
	github.com/asishrs/proxyless-grpc-lb/common v0.0.0-20200824005131-5718db7e19f5
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/envoyproxy/go-control-plane v0.9.6
	github.com/envoyproxy/protoc-gen-validate v0.4.1 // indirect
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/uuid v1.1.1
	github.com/googleapis/gnostic v0.4.0 // indirect
	github.com/json-iterator/go v1.1.10 // indirect
	github.com/kr/pretty v0.2.0 // indirect
	github.com/magiconair/properties v1.8.2 // indirect
	github.com/mitchellh/mapstructure v1.3.3 // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d // indirect
	golang.org/x/time v0.0.0-20200630173020-3af7569d3a1e // indirect
	google.golang.org/appengine v1.6.6 // indirect
	google.golang.org/genproto v0.0.0-20200825200019-8632dd797987 // indirect
	google.golang.org/grpc v1.31.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/ini.v1 v1.60.1 // indirect
	gopkg.in/yaml.v3 v3.0.0 // indirect
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
