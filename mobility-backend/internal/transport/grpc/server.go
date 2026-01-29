package grpc

import (
	"fmt"
	"net"

	"github.com/zeng/mobility-backend/internal/domain/fireball"
	"github.com/zeng/mobility-backend/internal/domain/location"
	"github.com/zeng/mobility-backend/internal/proto/gen"
	"github.com/zeng/mobility-backend/internal/transport/grpc/handlers"
	"github.com/zeng/mobility-backend/internal/transport/grpc/interceptors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ServerConfig holds config needed for gRPC server (avoids import cycle with app).
type ServerConfig struct {
	JWTSecret string
}

// Server wraps gRPC server and registers mobility services.
type Server struct {
	*grpc.Server
	log *zap.Logger
}

func NewServer(cfg ServerConfig, log *zap.Logger, locSvc *location.Service, fireballSvc *fireball.Service, pub fireball.LocationPublisher) *Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.AuthUnaryInterceptor(cfg.JWTSecret),
			interceptors.LoggingUnaryInterceptor(log),
			interceptors.RateLimitUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			interceptors.AuthStreamInterceptor(cfg.JWTSecret),
			interceptors.LoggingStreamInterceptor(log),
		),
	)
	locHandler := handlers.NewLocationHandler(locSvc, fireballSvc, pub)
	gen.RegisterLocationServiceServer(srv, locHandler)
	reflection.Register(srv)
	return &Server{Server: srv, log: log}
}

func (s *Server) Serve(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	return s.Server.Serve(lis)
}

func (s *Server) GracefulStop() {
	s.Server.GracefulStop()
}
