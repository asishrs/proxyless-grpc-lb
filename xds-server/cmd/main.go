package main

import (
	"context"
	"time"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	xds "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	config2 "github.com/grobza/proxyless-grpc-lb/hello-world/config"
	logger2 "github.com/grobza/proxyless-grpc-lb/hello-world/logger"
	"github.com/grobza/proxyless-grpc-lb/xds-server/internal/app"
	elogger "github.com/grobza/proxyless-grpc-lb/xds-server/internal/app"
	"go.uber.org/zap"
)

var (
	l elogger.Logger
)

func init() {
	l = elogger.Logger{}
}

func main() {

	config, err := config2.ReadConfig()
	if err != nil {
		logger2.Logger.Fatal("Unable to read config", zap.Error(err))
	}

	ctx := context.Background()

	logger2.Logger.Info("Starting control plane")
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

	logger2.Logger.Debug("Status", zap.Any("keys", snapshotCache.GetStatusKeys()))

	nodeID := config.GetString("nodeId")
	logger2.Logger.Info("Creating Node", zap.String("Id", nodeID))
	for {
		ss, err := app.GenerateSnapshot(config.GetStringSlice("upstreamServices"))
		if err != nil {
			logger2.Logger.Error("Error in Generating the SnapShot", zap.Error(err))
		} else {
			snapshotCache.SetSnapshot(nodeID, *ss)
			time.Sleep(60 * time.Second)
		}
	}

}
