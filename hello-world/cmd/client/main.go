package main

import (
	"fmt"

	"google.golang.org/grpc"

	"go.uber.org/zap"

	config "github.com/grobza/proxyless-grpc-lb/common/pkg/config"
	logger "github.com/grobza/proxyless-grpc-lb/common/pkg/logger"
	client "github.com/grobza/proxyless-grpc-lb/hello-world/internal/app/client"
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
