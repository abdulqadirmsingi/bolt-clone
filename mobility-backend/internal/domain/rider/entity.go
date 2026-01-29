package rider

// Entity represents a rider (persisted in PostgreSQL).
// Live session/presence can be in Redis; no GPS stream written to DB.
type Entity struct {
	ID     string
	UserID string
}
