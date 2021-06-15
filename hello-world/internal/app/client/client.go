package client

import (
	"os"
	"os/signal"
	"syscall"

	"time"

	"go.uber.org/zap"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	logger "github.com/grobza/proxyless-grpc-lb/common/pkg/logger"
	hello "github.com/grobza/proxyless-grpc-lb/hello-world/internal/app/http/rpc"
	helper "github.com/grobza/proxyless-grpc-lb/hello-world/internal/pkg"
	_ "google.golang.org/grpc/xds"
)


var (
	conn *grpc.ClientConn
)

func StartClient(server string) {
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	logger.Logger.Info("Starting gRPC Client", zap.String("host", server))

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(server, opts...)

	if err != nil {
		logger.Logger.Fatal("Unable to start", zap.Error(err))
	}

	go helper.ShutdownClient(stop, conn)

	c := hello.NewHelloClient(conn)
	ctx := context.Background()
	for {
		msg, err := c.SayHello(ctx, &hello.HelloRequest{Name: "gRPC Proxyless LB"})
		if err != nil {
			logger.Logger.Error("Unable to send Hello message", zap.Any("Response", msg), zap.Error(err))
		} else {
			logger.Logger.Info("Hello Response", zap.Any("Response", msg))
		}
		time.Sleep(3 * time.Second)
	}
}
