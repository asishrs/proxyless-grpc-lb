package main

import (
	"github.com/asishrs/proxyless-grpc-lb/common/pkg/config"
	"github.com/asishrs/proxyless-grpc-lb/common/pkg/logger"
	"github.com/asishrs/proxyless-grpc-lb/hello-world/internal/app/server"
	"go.uber.org/zap"
)

func main() {

	config, err := config.ReadConfig()
	if err != nil {
		logger.Logger.Fatal("Unable to read config", zap.Error(err))
	}

	server.StartServer(config.GetInt("port"))
}
