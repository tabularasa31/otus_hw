package grpcserver

import (
	"context"

	"go.uber.org/zap"
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
func (im *InterceptorManager) Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	md, _ := metadata.FromIncomingContext(ctx)
	reply, err := handler(ctx, req)
	im.logger.Infof("Server: %v Method: %s, Metadata: %v, Err: %v", info.Server, info.FullMethod, md, err)

	return reply, err
}
