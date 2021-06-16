package helper

import (
	"context"

	logger2 "github.com/grobza/proxyless-grpc-lb/hello-world/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

// Health struct
type Health struct{}

// Check is for HealthCheck
func (h *Health) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	logger2.Logger.Debug("health check", zap.Any("status", grpc_health_v1.HealthCheckResponse_SERVING))
	return &grpc_health_v1.HealthCheckResponse{
		Status: grpc_health_v1.HealthCheckResponse_SERVING,
	}, nil
}

// Watch is used by clients to get updates
func (h *Health) Watch(req *grpc_health_v1.HealthCheckRequest, w grpc_health_v1.Health_WatchServer) error {
	return status.Error(codes.Unimplemented, "Not Supported")
}
