package app

import (
	"context"
	"os"
	"sync"

	"github.com/zeng/mobility-backend/internal/domain/fireball"
	"github.com/zeng/mobility-backend/internal/domain/location"
	"github.com/zeng/mobility-backend/internal/infrastructure/postgres"
	"github.com/zeng/mobility-backend/internal/infrastructure/redis"
	"github.com/zeng/mobility-backend/internal/transport/grpc"
	"go.uber.org/zap"
)

type Container struct {
	cfg      *Config
	log      *zap.Logger
	pg       *postgres.DB
	rdb      *redis.Client
	locSvc   *location.Service
	fireball *fireball.Service
	grpcSrv  *grpc.Server
}

func NewContainer(ctx context.Context) (*Container, error) {
	cfgPath := "configs/local.yaml"
	if p := osGetenv("CONFIG_PATH"); p != "" {
		cfgPath = p
	}
	cfg, err := LoadConfig(cfgPath)
	if err != nil {
		return nil, err
	}
	log, _ := zap.NewProduction()
	if cfg.Env == "local" || cfg.Env == "dev" {
		log, _ = zap.NewDevelopment()
	}

	pg, err := postgres.Connect(ctx, cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	rdb, err := redis.Connect(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		return nil, err
	}

	// H3 + geo
	geoRepo := redis.NewGeoRepo(rdb)
	h3Index := location.NewH3Index(cfg.H3.Resolution, cfg.H3.DefaultKRing)
	locSvc := location.NewService(h3Index, geoRepo)

	// Fireball: threshold logic + publisher
	threshold := fireball.NewThreshold(
		cfg.Fireball.MinDistanceMeters,
		cfg.Fireball.MinHeadingDegrees,
		cfg.Fireball.MaxSilenceSeconds,
		cfg.Fireball.ThrottleMs,
	)
	publisher := redis.NewFireballPublisher(rdb)
	fireballSvc := fireball.NewService(threshold, locSvc, publisher)

	grpcSrv := grpc.NewServer(grpc.ServerConfig{JWTSecret: cfg.JWT.Secret}, log, locSvc, fireballSvc, publisher)

	return &Container{
		cfg:      cfg,
		log:      log,
		pg:       pg,
		rdb:      rdb,
		locSvc:   locSvc,
		fireball: fireballSvc,
		grpcSrv:  grpcSrv,
	}, nil
}

func (c *Container) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := c.grpcSrv.Serve(c.cfg.GRPC.Port); err != nil {
			c.log.Error("grpc serve", zap.Error(err))
		}
	}()
	c.log.Info("mobility-backend started", zap.Int("grpc_port", c.cfg.GRPC.Port))
	wg.Wait()
}

func (c *Container) Shutdown(ctx context.Context) {
	c.log.Info("shutting down")
	c.grpcSrv.GracefulStop()
	_ = c.rdb.Close()
	_ = c.pg.Close()
}

// osGetenv for testability
var osGetenv = os.Getenv
