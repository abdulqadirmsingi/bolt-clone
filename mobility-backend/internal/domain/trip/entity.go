package trip

// Entity represents a trip (PostgreSQL).
type Entity struct {
	ID              string
	RiderID         string
	DriverID        string
	Status          string // requested, accepted, in_progress, completed, cancelled
	PickupLat       float64
	PickupLng       float64
	DropoffLat      float64
	DropoffLng      float64
	CurrentStopIndex int
	OrderedStops   []Stop
	ETASeconds     int
}

type Stop struct {
	Lat      float64
	Lng      float64
	Label    string
	Sequence int
}
