package main

import (
	"fmt"

	"github.com/asishrs/proxyless-grpc-lb/common/pkg/config"
	"github.com/asishrs/proxyless-grpc-lb/common/pkg/logger"
	"github.com/asishrs/proxyless-grpc-lb/hello-world/internal/app/client"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	conn *grpc.ClientConn
)

func main() {

	config, err := config.ReadConfig()
	if err != nil {
		logger.Logger.Fatal("Unable to read config", zap.Error(err))
	}

	client.StartClient(fmt.Sprintf("%s", config.GetString("hello.host")))
}
