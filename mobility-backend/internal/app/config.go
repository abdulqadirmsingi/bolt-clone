package app

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Env      string       `yaml:"env"`
	GRPC     GRPCConfig   `yaml:"grpc"`
	HTTP     HTTPConfig   `yaml:"http"`
	Postgres PostgresConfig `yaml:"postgres"`
	Redis    RedisConfig  `yaml:"redis"`
	Fireball FireballConfig `yaml:"fireball"`
	H3       H3Config     `yaml:"h3"`
	JWT      JWTConfig    `yaml:"jwt"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type HTTPConfig struct {
	Port int `yaml:"port"`
}

type PostgresConfig struct {
	DSN string `yaml:"dsn"`
}

type RedisConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type FireballConfig struct {
	// MinDistanceMeters: push only when driver moved at least this much
	MinDistanceMeters float64 `yaml:"min_distance_meters"`
	// MinHeadingDegrees: push only when heading changed by this much
	MinHeadingDegrees float64 `yaml:"min_heading_degrees"`
	// MaxSilenceSeconds: force push if no update for this long (keep-alive)
	MaxSilenceSeconds int `yaml:"max_silence_seconds"`
	// ThrottleMs: minimum interval between pushes per driver
	ThrottleMs int `yaml:"throttle_ms"`
}

type H3Config struct {
	// Resolution 9 ≈ 0.1 km² cells; 10 ≈ 0.03 km². Use 9 for driver discovery.
	Resolution int `yaml:"resolution"`
	// DefaultKRing: number of rings for "nearby drivers" (K=2 → 19 cells)
	DefaultKRing int `yaml:"default_k_ring"`
}

type JWTConfig struct {
	Secret     string `yaml:"secret"`
	ExpireMins int    `yaml:"expire_mins"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}
