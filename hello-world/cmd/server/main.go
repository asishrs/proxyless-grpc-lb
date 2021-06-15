package main

import (
	config "github.com/grobza/proxyless-grpc-lb/common/pkg/config"
	logger "github.com/grobza/proxyless-grpc-lb/common/pkg/logger"
	server "github.com/grobza/proxyless-grpc-lb/hello-world/internal/app/server"
	"go.uber.org/zap"
)

func main() {

	config, err := config.ReadConfig()
	if err != nil {
		logger.Logger.Fatal("Unable to read config", zap.Error(err))
	}

	server.StartServer(config.GetInt("port"))
}
