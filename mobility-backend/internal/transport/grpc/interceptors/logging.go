package interceptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingUnaryInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Debug("grpc unary", zap.String("method", info.FullMethod))
		return handler(ctx, req)
	}
}

func LoggingStreamInterceptor(log *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		log.Debug("grpc stream", zap.String("method", info.FullMethod))
		return handler(srv, ss)
	}
}
