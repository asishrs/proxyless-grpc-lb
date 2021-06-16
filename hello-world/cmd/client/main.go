package main

import (
	"fmt"

	config2 "github.com/grobza/proxyless-grpc-lb/hello-world/config"
	"github.com/grobza/proxyless-grpc-lb/hello-world/internal/app/client"
	logger2 "github.com/grobza/proxyless-grpc-lb/hello-world/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	conn *grpc.ClientConn
)

func main() {

	config, err := config2.ReadConfig()
	if err != nil {
		logger2.Logger.Fatal("Unable to read config", zap.Error(err))
	}

	client.StartClient(fmt.Sprintf("%s", config.GetString("hello.host")))
}
