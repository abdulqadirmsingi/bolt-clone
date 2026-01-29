package app

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/zeng/mobility-backend/internal/domain/driver"
	"github.com/zeng/mobility-backend/internal/domain/fireball"
	"github.com/zeng/mobility-backend/internal/domain/location"
	"github.com/zeng/mobility-backend/internal/domain/trip"
	"github.com/zeng/mobility-backend/internal/infrastructure/postgres"
	"github.com/zeng/mobility-backend/internal/infrastructure/redis"
	"github.com/zeng/mobility-backend/internal/transport/grpc"
	"go.uber.org/zap"
)

type Container struct {
	cfg         *Config
	log         *zap.Logger
	pg          *postgres.DB
	rdb         *redis.Client
	locSvc      *location.Service
	fireball    *fireball.Service
	availSvc    *driver.AvailabilityService
	heartbeat   *HeartbeatChecker
	grpcSrv     *grpc.Server
	cancelRun   context.CancelFunc
}

func NewContainer(ctx context.Context) (*Container, error) {
	loadEnvIfExists()
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

	pgDSN := postgresDSN(cfg.Postgres.DSN)
	pg, err := postgres.Connect(ctx, pgDSN)
	if err != nil {
		return nil, err
	}
	if err := pg.RunMigrations(ctx); err != nil {
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

	// Driver presence & availability
	availRepo := redis.NewAvailabilityRepo(rdb)
	presenceNotifier := &redis.NoopPresenceNotifier{}
	availSvc := driver.NewAvailabilityService(availRepo, presenceNotifier, locSvc)
	fireballSvc.SetAvailabilityChecker(newAvailabilityCheckerAdapter(availSvc))

	// Heartbeat checker: mark drivers OFFLINE when no heartbeat within timeout
	timeoutSec := cfg.Driver.HeartbeatTimeoutSeconds
	if timeoutSec <= 0 {
		timeoutSec = 45
	}
	intervalSec := timeoutSec / 2
	if intervalSec < 15 {
		intervalSec = 15
	}
	heartbeatChecker := NewHeartbeatChecker(availRepo, availSvc, timeoutSec, intervalSec)

	// Trip: Postgres repo (migrations create trips/trip_stops). Use for local and production.
	tripRepo := postgres.NewTripRepo(pg)
	tripSvc := trip.NewService(tripRepo)

	grpcSrv := grpc.NewServer(grpc.ServerConfig{JWTSecret: cfg.JWT.Secret}, log, locSvc, fireballSvc, publisher, tripSvc, availSvc)

	return &Container{
		cfg:       cfg,
		log:       log,
		pg:        pg,
		rdb:       rdb,
		locSvc:    locSvc,
		fireball:  fireballSvc,
		availSvc:  availSvc,
		heartbeat: heartbeatChecker,
		grpcSrv:   grpcSrv,
	}, nil
}

func (c *Container) Start(ctx context.Context) {
	runCtx, cancel := context.WithCancel(ctx)
	c.cancelRun = cancel
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		c.heartbeat.Run(runCtx)
	}()
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
	if c.cancelRun != nil {
		c.cancelRun()
	}
	c.grpcSrv.GracefulStop()
	_ = c.rdb.Close()
	c.pg.Close()
}

// osGetenv for testability
var osGetenv = os.Getenv

// postgresDSN returns Postgres DSN: POSTGRES_DSN env, else DB_* env vars (Django-style), else config.
func postgresDSN(defaultDSN string) string {
	if dsn := osGetenv("POSTGRES_DSN"); dsn != "" {
		return dsn
	}
	user := osGetenv("DB_USER")
	host := osGetenv("DB_HOST")
	if user == "" && host == "" {
		return defaultDSN
	}
	name := osGetenv("DB_NAME")
	if name == "" {
		name = "mobility"
	}
	if host == "" {
		host = "localhost"
	}
	port := osGetenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	password := osGetenv("DB_PASSWORD")
	if password != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, name)
	}
	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", user, host, port, name)
}
