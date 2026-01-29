package interceptors

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthUnaryInterceptor(jwtSecret string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if jwtSecret == "" {
			return handler(ctx, req)
		}
		ctx, err := extractAndValidateJWT(ctx, jwtSecret)
		if err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}

func AuthStreamInterceptor(jwtSecret string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		if jwtSecret == "" {
			return handler(srv, ss)
		}
		ctx, err := extractAndValidateJWT(ss.Context(), jwtSecret)
		if err != nil {
			return err
		}
		return handler(srv, &streamWithContext{ServerStream: ss, ctx: ctx})
	}
}

func extractAndValidateJWT(ctx context.Context, secret string) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}
	vals := md.Get("authorization")
	if len(vals) == 0 {
		return nil, status.Error(codes.Unauthenticated, "missing authorization")
	}
	token := strings.TrimPrefix(vals[0], "Bearer ")
	if token == vals[0] {
		return nil, status.Error(codes.Unauthenticated, "invalid authorization format")
	}
	// TODO: validate JWT with jwtSecret and set subject in context
	_ = token
	return ctx, nil
}

type streamWithContext struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *streamWithContext) Context() context.Context {
	return s.ctx
}
