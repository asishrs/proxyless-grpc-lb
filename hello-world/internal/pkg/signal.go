package helper

import (
	"os"
	"time"

	logger2 "github.com/grobza/proxyless-grpc-lb/hello-world/logger"
	"google.golang.org/grpc"
)

// ShutdownServer gracefully
func ShutdownServer(stop chan os.Signal, server *grpc.Server) {
	<-stop
	logger2.Logger.Debug("Stopping Server")
	time.Sleep(60 * time.Second)
	server.GracefulStop()
}

// ShutdownClient gracefully
func ShutdownClient(stop chan os.Signal, connection *grpc.ClientConn) {
	<-stop
	logger2.Logger.Debug("Stopping Client")
	time.Sleep(60 * time.Second)
	err := connection.Close()
	if err != nil {
		panic(err)
	}
}
