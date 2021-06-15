package main

import (
	"context"
	"time"

	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v2"

	app "github.com/grobza/proxyless-grpc-lb/xds-server/internal/app"

	"go.uber.org/zap"

	config "github.com/grobza/proxyless-grpc-lb/common/pkg/config"
	logger "github.com/grobza/proxyless-grpc-lb/common/pkg/logger"

	elogger "github.com/grobza/proxyless-grpc-lb/xds-server/internal/app"
)

var (
	l elogger.Logger
)

func init() {
	l = elogger.Logger{}
}

func main() {

	config, err := config.ReadConfig()
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
