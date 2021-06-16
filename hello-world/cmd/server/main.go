package main

import (
	config2 "github.com/grobza/proxyless-grpc-lb/hello-world/config"
	"github.com/grobza/proxyless-grpc-lb/hello-world/internal/app/server"
	logger2 "github.com/grobza/proxyless-grpc-lb/hello-world/logger"
	"go.uber.org/zap"
)

func main() {

	config, err := config2.ReadConfig()
	if err != nil {
		logger2.Logger.Fatal("Unable to read config", zap.Error(err))
	}

	server.StartServer(config.GetInt("port"))
}
