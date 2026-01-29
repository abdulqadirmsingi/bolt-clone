package driver

// Entity represents a driver (persisted in PostgreSQL; live location in Redis).
type Entity struct {
	ID        string
	UserID    string
	Status    string // offline, available, on_trip
	VehicleID string
}
