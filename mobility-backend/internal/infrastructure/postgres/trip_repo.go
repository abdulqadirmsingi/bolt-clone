package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/zeng/mobility-backend/internal/domain/trip"
)

// TripRepo implements trip.Repository using PostgreSQL. Use for production.
type TripRepo struct {
	db *DB
}

func NewTripRepo(db *DB) *TripRepo {
	return &TripRepo{db: db}
}

func (r *TripRepo) Create(ctx context.Context, e *trip.Entity) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO trips (id, rider_id, driver_id, status, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, current_stop_index, eta_seconds)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`, e.ID, e.RiderID, e.DriverID, e.Status, e.PickupLat, e.PickupLng, e.DropoffLat, e.DropoffLng, e.CurrentStopIndex, e.ETASeconds)
	if err != nil {
		return err
	}
	for _, s := range e.OrderedStops {
		_, err = r.db.Exec(ctx, `INSERT INTO trip_stops (trip_id, sequence, lat, lng, label) VALUES ($1, $2, $3, $4, $5)`,
			e.ID, s.Sequence, s.Lat, s.Lng, s.Label)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TripRepo) GetByID(ctx context.Context, id string) (*trip.Entity, error) {
	var e trip.Entity
	err := r.db.QueryRow(ctx, `
		SELECT id, rider_id, driver_id, status, pickup_lat, pickup_lng, dropoff_lat, dropoff_lng, current_stop_index, eta_seconds
		FROM trips WHERE id = $1
	`, id).Scan(&e.ID, &e.RiderID, &e.DriverID, &e.Status, &e.PickupLat, &e.PickupLng, &e.DropoffLat, &e.DropoffLng, &e.CurrentStopIndex, &e.ETASeconds)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	rows, err := r.db.Query(ctx, `SELECT sequence, lat, lng, label FROM trip_stops WHERE trip_id = $1 ORDER BY sequence`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var s trip.Stop
		if err := rows.Scan(&s.Sequence, &s.Lat, &s.Lng, &s.Label); err != nil {
			return nil, err
		}
		e.OrderedStops = append(e.OrderedStops, s)
	}
	return &e, nil
}

func (r *TripRepo) UpdateStatus(ctx context.Context, id, status string) error {
	_, err := r.db.Exec(ctx, `UPDATE trips SET status = $1, updated_at = now() WHERE id = $2`, status, id)
	return err
}

func (r *TripRepo) UpdateCurrentStop(ctx context.Context, id string, currentStopIndex int) error {
	_, err := r.db.Exec(ctx, `UPDATE trips SET current_stop_index = $1, updated_at = now() WHERE id = $2`, currentStopIndex, id)
	return err
}

func (r *TripRepo) UpdateETA(ctx context.Context, id string, etaSeconds int) error {
	_, err := r.db.Exec(ctx, `UPDATE trips SET eta_seconds = $1, updated_at = now() WHERE id = $2`, etaSeconds, id)
	return err
}

var _ trip.Repository = (*TripRepo)(nil)
