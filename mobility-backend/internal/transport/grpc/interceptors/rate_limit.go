package interceptors

import (
	"context"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RateLimitUnaryInterceptor applies a simple per-connection rate limit (stub).
func RateLimitUnaryInterceptor() grpc.UnaryServerInterceptor {
	var mu sync.Mutex
	allowed := make(map[string]int)
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: key by peer/subject, enforce limit (e.g. 100 req/s per client)
		mu.Lock()
		allowed[""]++
		mu.Unlock()
		return handler(ctx, req)
	}
}

