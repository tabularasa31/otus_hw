package grpcserver

import (
	"context"
	"go.uber.org/zap"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// InterceptorManager .
type InterceptorManager struct {
	logger zap.SugaredLogger
}

// NewInterceptorManager - InterceptorManager constructor .
func NewInterceptorManager(logger zap.SugaredLogger) *InterceptorManager {
	return &InterceptorManager{logger: logger}
}

// Logger Interceptor .
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	start := time.Now()
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Info("Method: %s, Time: %v, Metadata: %v, Err: %v", info.FullMethod, time.Since(start), md, err)

	return reply, err
}
