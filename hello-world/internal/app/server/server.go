package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"

	logger "github.com/grobza/proxyless-grpc-lb/common/pkg/logger"
	hello "github.com/grobza/proxyless-grpc-lb/hello-world/internal/app/http/rpc"
	helper "github.com/grobza/proxyless-grpc-lb/hello-world/internal/pkg"
)

type server struct {
}

func (s *server) SayHello(ctx context.Context, in *hello.HelloRequest) (*hello.HelloResponse, error) {
	logger.Logger.Debug("Incoming Request", zap.String("name", in.Name))

	var host, err = os.Hostname()
	if err != nil {
		logger.Logger.Debug("Unable to get hostname", zap.Error(err))
	}
	return &hello.HelloResponse{Message: fmt.Sprintf("Hello %s from host %s", in.Name, host)}, nil
}

// StartServer starts the gRPC server
func StartServer(port int) {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	logger.Logger.Info("Starting gRPC Server", zap.Int("port", port))

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		logger.Logger.Info("failed to listen", zap.Error(err))
	}

	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	hello.RegisterHelloServer(grpcServer, &server{})
	grpc_health_v1.RegisterHealthServer(grpcServer, &helper.Health{})
	go helper.ShutdownServer(stop, grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Logger.Fatal("Failed to start gRPC Server")
	}

}
