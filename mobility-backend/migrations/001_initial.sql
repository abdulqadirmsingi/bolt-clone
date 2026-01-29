-- Initial schema for mobility. Run this in DBeaver (or the backend runs it on startup).
-- Live driver location and availability are in Redis; this is for driver/trip identity.

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
    rider_id            TEXT NOT NULL REFERENCES riders(id),
    driver_id           TEXT NOT NULL REFERENCES drivers(id),
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
    trip_id     TEXT NOT NULL REFERENCES trips(id) ON DELETE CASCADE,
    sequence    INT NOT NULL,
    lat         DOUBLE PRECISION NOT NULL,
    lng         DOUBLE PRECISION NOT NULL,
    label       TEXT NOT NULL DEFAULT '',
    PRIMARY KEY (trip_id, sequence)
);

CREATE INDEX IF NOT EXISTS idx_trips_rider_id ON trips(rider_id);
CREATE INDEX IF NOT EXISTS idx_trips_driver_id ON trips(driver_id);
CREATE INDEX IF NOT EXISTS idx_trips_status ON trips(status);
