package main

import (
	"context"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/grobza/proxyless-grpc-lb/xds-server/internal/app"
	elogger "github.com/grobza/proxyless-grpc-lb/xds-server/internal/app"
	"github.com/grobza/proxyless-grpc-lb/xds-server/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	l elogger.Logger
)

func init() {
	l = elogger.Logger{}
}

// ReadConfig reads the config data from file
func ReadConfig() (*viper.Viper, error) {
	logger.Logger.Debug("Reading configuration", zap.String("file", "/var/run/config/app.yaml"))
	v := viper.New()
	v.SetConfigFile("/var/run/config/app.yaml")
	v.AutomaticEnv()
	err := v.ReadInConfig()
	return v, err

}

func main() {

	config, err := ReadConfig()
	if err != nil {
		logger.Logger.Fatal("Unable to read config", zap.Error(err))
	}

	ctx := context.Background()

	logger.Logger.Info("Starting control plane")
	signal := make(chan struct{})
	cb := &app.Callbacks{
		Signal:   signal,
		Fetches:  0,
		Requests: 0,
	}

	snapshotCache := cache.NewSnapshotCache(true, cache.IDHash{}, l)

	srv := xds.NewServer(ctx, snapshotCache, cb)

	go app.RunManagementServer(ctx, srv, uint(config.GetInt("managementServer.port")), uint32(config.GetInt("maxConcurrentStreams")))
	<-signal

	cb.Report()

	logger.Logger.Debug("Status", zap.Any("keys", snapshotCache.GetStatusKeys()))

	nodeID := config.GetString("nodeId")
	logger.Logger.Info("Creating Node", zap.String("Id", nodeID))
	for {
		ss, err := app.GenerateSnapshot(config.GetStringSlice("upstreamServices"))
		if err != nil {
			logger.Logger.Error("Error in Generating the SnapShot", zap.Error(err))
		} else {
			snapshotCache.SetSnapshot(nodeID, *ss)
			time.Sleep(60 * time.Second)
		}
	}

}
