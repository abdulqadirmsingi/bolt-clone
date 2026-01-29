package postgres

import (
	"context"
	"strings"
)

// initialSQL creates drivers, riders, trips, trip_stops. Live location/availability are in Redis.
const initialSQL = `
CREATE TABLE IF NOT EXISTS drivers (
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL,
    vehicle_id  TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS riders (
    id          TEXT PRIMARY KEY,
    user_id     TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS trips (
    id                  TEXT PRIMARY KEY,
    rider_id            TEXT NOT NULL,
    driver_id           TEXT NOT NULL,
    status              TEXT NOT NULL DEFAULT 'requested',
    pickup_lat          DOUBLE PRECISION NOT NULL,
    pickup_lng          DOUBLE PRECISION NOT NULL,
    dropoff_lat         DOUBLE PRECISION NOT NULL,
    dropoff_lng         DOUBLE PRECISION NOT NULL,
    current_stop_index  INT NOT NULL DEFAULT 0,
    eta_seconds         INT NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS trip_stops (
    trip_id     TEXT NOT NULL,
    sequence    INT NOT NULL,
    lat         DOUBLE PRECISION NOT NULL,
    lng         DOUBLE PRECISION NOT NULL,
    label       TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (trip_id, sequence)
);
CREATE INDEX IF NOT EXISTS idx_trips_rider_id ON trips(rider_id);
CREATE INDEX IF NOT EXISTS idx_trips_driver_id ON trips(driver_id);
CREATE INDEX IF NOT EXISTS idx_trips_status ON trips(status);
`

// RunMigrations runs the initial schema (idempotent). Call once after Connect.
func (d *DB) RunMigrations(ctx context.Context) error {
	for _, stmt := range strings.Split(strings.TrimSpace(initialSQL), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := d.Exec(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}
